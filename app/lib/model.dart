import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:scoped_model/scoped_model.dart';

class Entity {
  final String text;
  final String type;
  final double score;

  Entity({
    required this.text,
    required this.type,
    required this.score,
  });

  factory Entity.fromJson(Map<String, dynamic> json) {
    return Entity(
      text: json['text'] as String,
      type: json['type'] as String,
      score: json['score'] as double,
    );
  }
}

class KeyTerms {
  final List<Entity> entities;

  KeyTerms({
    required this.entities,
  });

  factory KeyTerms.fromJson(Map<String, dynamic> json) {
    return KeyTerms(
      entities: (json['entities'] as List<dynamic>)
          .map((e) => Entity.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }
}

class ReferenceMaterial {
  final String matched;
  final List<String> terms;
  final List<String> images;

  ReferenceMaterial(this.matched, this.terms, this.images);

  factory ReferenceMaterial.fromJson(Map<String, dynamic> json) {
    return ReferenceMaterial(
      json['matched'] as String,
      (json['terms'] as List<dynamic>).cast<String>(),
      (json['images'] as List<dynamic>).cast<String>(),
    );
  }
}

class SearchImage {
  final String contentUrl;
  final String contentSize;
  final String thumbnailUrl;
  final String hostPageUrl;
  final String encodingFormat;
  final int width;
  final int height;
  final String accentColor;
  bool? isLiked;

  SearchImage({
    required this.contentUrl,
    required this.contentSize,
    required this.thumbnailUrl,
    required this.hostPageUrl,
    required this.encodingFormat,
    required this.width,
    required this.height,
    required this.accentColor,
    this.isLiked,
  });

  factory SearchImage.fromJson(Map<String, dynamic> json) {
    return SearchImage(
      contentUrl: json['contentURL'] as String,
      contentSize: json['contentSize'] as String,
      thumbnailUrl: json['thumbnailURL'] as String,
      hostPageUrl: json['hostPageURL'] as String,
      encodingFormat: json['encodingFormat'] as String,
      width: json['width'] as int,
      height: json['height'] as int,
      accentColor: json['accentColor'] as String,
      isLiked: json['isLiked'] as bool?,
    );
  }
}

void checkHttpStatus(Response response) {
  if (response.statusCode != 200 && response.statusCode != 202) {
    throw ErrorSummary("status ${response.statusCode}: ${response.body}");
  }
}

Future<List<SearchImage>> search(String query) async {
  final url = Uri.https('ts.beebs.dev', '/api/search', {
    'q': query,
  });
  final response = await http.get(url);
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as List? ?? [];
  return decodedResponse.map((e) => SearchImage.fromJson(e)).toList();
}

class MyModel extends Model {
  bool _isConnected = false;
  WebSocketChannel? _channel;
  String? _transcript = "";
  final List<ReferenceMaterial> _referenceMaterials = [];
  KeyTerms? _keyTerms;
  Entity? _selectedEntity;
  List<SearchImage>? _searchImages;

  bool get isConnected => _isConnected;
  String? get transcript => _transcript;
  List<ReferenceMaterial> get referenceMaterials => _referenceMaterials;
  KeyTerms? get keyTerms => _keyTerms;
  Entity? get selectedEntity => _selectedEntity;
  List<SearchImage>? get searchImages => _searchImages;

  Future<void> searchForEntity(Entity entity) async {
    _searchImages = await search(entity.text);
    notifyListeners();
  }

  set selectedEntity(Entity? entity) {
    _selectedEntity = entity;
    if (entity != null) {
      searchForEntity(entity);
    }
    notifyListeners();
  }

  MyModel() {
    connectWebSock();
  }

  void onConnect() {
    _isConnected = true;
    notifyListeners();
  }

  void onDisconnect() {
    _isConnected = false;
    notifyListeners();
  }

  void likeImage(SearchImage img, bool isLiked) {
    img.isLiked = isLiked;
    notifyListeners();
  }

  Future<void> connectWebSock() async {
    onDisconnect();
    _channel?.sink.close();
    _channel = WebSocketChannel.connect(
      Uri.parse('wss://ts.beebs.dev/ws'),
    );
    _channel!.stream.listen(
      (message) => handleWebSockMessage(message),
      onError: (err) async {
        print('websock error: $err');
        displayTranscript('websocket error: $err');
        onDisconnect();
        await Future.delayed(const Duration(seconds: 2));
        connectWebSock();
      },
      onDone: () async {
        onDisconnect();
        await Future.delayed(const Duration(seconds: 2));
        connectWebSock();
      },
    );
  }

  void handleWebSockMessage(dynamic message) {
    if (_channel == null) return;
    if (!_isConnected) onConnect();
    final obj = jsonDecode(message) as Map<String, dynamic>;
    switch (obj['type']) {
      case 'ping':
        _channel!.sink.add(jsonEncode({'type': 'pong'}));
        break;
      case 'transcript':
        // received transcript
        displayTranscript(obj['payload']['text'] as String);
        break;
      case 'ref':
        // received reference material
        displayReference(ReferenceMaterial.fromJson(obj['payload']));
        break;
      case 'keyterms':
        // received key terms
        displayKeyTerms(KeyTerms.fromJson(obj['payload']));
        break;
      default:
        break;
    }
  }

  void displayTranscript(String transcript) {
    _transcript = transcript;
    notifyListeners();
  }

  void displayReference(ReferenceMaterial ref) {
    // display reference material
    _referenceMaterials.add(ref);
    while (_referenceMaterials.length > 15) {
      // limit the number of reference material displayed
      _referenceMaterials.removeAt(0);
    }
    notifyListeners();
  }

  void displayKeyTerms(KeyTerms keyTerms) {
    _keyTerms = keyTerms;
    notifyListeners();
  }
}
