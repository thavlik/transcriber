import 'dart:convert';
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

class MyModel extends Model {
  bool _isConnected = false;
  WebSocketChannel? _channel;
  String? _transcript = "";
  final List<ReferenceMaterial> _referenceMaterials = [];
  KeyTerms? _keyTerms;

  /*
    ReferenceMaterial("vertebral arch", [
      "vertebral arch",
    ], [
      "https://refmat.nyc3.digitaloceanspaces.com/lumbar-vertebra-vertebral-arch-superior-view-745x550.png",
      "https://refmat.nyc3.digitaloceanspaces.com/General-Structure-of-a-Vertebrae.jpg",
    ]),
    ReferenceMaterial("ligamentum flavum", [
      "ligamentum flavum",
      "ligamentum",
      "flavum",
    ], [
      "https://refmat.nyc3.digitaloceanspaces.com/ligamentum-flavum-1024x670.jpg",
      "https://refmat.nyc3.digitaloceanspaces.com/LigamentumFlavum.png",
    ]),
    ReferenceMaterial("facet joint", [
      "facet joint",
    ], [
      "https://refmat.nyc3.digitaloceanspaces.com/facet_joints_related_spine_structures_shutterstock_157672247.jpg",
      "https://refmat.nyc3.digitaloceanspaces.com/Thoracic-Facet-Syndrome.jpg",
    ]),
  ];
  */

  bool get isConnected => _isConnected;
  String? get transcript => _transcript;
  List<ReferenceMaterial> get referenceMaterials => _referenceMaterials;
  KeyTerms? get keyTerms => _keyTerms;

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
