import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:scoped_model/scoped_model.dart';
import 'api.dart' as api;

class MyModel extends Model {
  bool _isConnected = false;
  WebSocketChannel? _channel;
  String _transcript = "";
  final List<api.ReferenceMaterial> _referenceMaterials = [];
  api.KeyTerms? _keyTerms;
  api.Entity? _selectedEntity;
  List<api.SearchImage>? _searchImages;
  List<api.SearchImage>? _radiologySearchImages;
  final Map<String, String> _definitions = {};
  final Set<String> _fetchingTerms = {};

  bool get isConnected => _isConnected;
  String get transcript => _transcript;
  List<api.ReferenceMaterial> get referenceMaterials => _referenceMaterials;
  api.KeyTerms? get keyTerms => _keyTerms;
  api.Entity? get selectedEntity => _selectedEntity;
  List<api.SearchImage>? get searchImages => _searchImages;
  List<api.SearchImage>? get radiologySearchImages => _radiologySearchImages;

  Future<void> search(String text) async {
    _searchImages = await api.search(text);
    notifyListeners();
  }

  Future<void> searchRadiology(String text) async {
    _radiologySearchImages = await api.search('$text radiology');
    notifyListeners();
  }

  set selectedEntity(api.Entity? entity) {
    _selectedEntity = entity;
    if (entity != null) {
      search(entity.text);
      if (entity.type == 'SYSTEM_ORGAN_SITE') searchRadiology(entity.text);
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

  String? getDefinition(String term) {
    final def = _definitions[term];
    if (def == null) define(term);
    return def;
  }

  bool termIsFetching(String term) => _fetchingTerms.contains(term);

  Future<void> define(String term) async {
    if (_definitions.containsKey(term) || _fetchingTerms.contains(term)) return;
    _fetchingTerms.add(term);
    final def = await api.define(term);
    _fetchingTerms.remove(term);
    _definitions[term] = def;
    notifyListeners();
  }

  void likeImage(api.SearchImage img, bool isLiked) async {
    //await api.likeImage(img, isLiked);
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
        displayReference(api.ReferenceMaterial.fromJson(obj['payload']));
        break;
      case 'keyterms':
        // received key terms
        displayKeyTerms(api.KeyTerms.fromJson(obj['payload']));
        break;
      default:
        break;
    }
  }

  void displayTranscript(String transcript) {
    _transcript = transcript;
    notifyListeners();
  }

  void displayReference(api.ReferenceMaterial ref) {
    // display reference material
    _referenceMaterials.add(ref);
    while (_referenceMaterials.length > 15) {
      // limit the number of reference material displayed
      _referenceMaterials.removeAt(0);
    }
    notifyListeners();
  }

  void displayKeyTerms(api.KeyTerms keyTerms) {
    _keyTerms ??= api.KeyTerms(entities: []);
    _keyTerms!.integrate(keyTerms.entities);
    _keyTerms!.prune();
    notifyListeners();
  }
}
