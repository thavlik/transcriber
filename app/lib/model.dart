import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:scoped_model/scoped_model.dart';
import 'api.dart' as api;

const exampleTranscript =
    "You want to aggressively attack the osteophytes to the reduce the diameter of the humeral head. All of that will give you a better understanding of the true anatomy as well as give you better glenoid exposure, which is the key in all these shoulder replacements - it's getting better access to the glenoid. So this crown is what we call a crown or ring of osteophytes,";

const exampleDefinition =
    'The glenoid is the shallow, concave cavity or socket found in the shoulder blade (scapula). It is part of the glenohumeral joint, where it articulates with the head of the humerus, forming the shoulder joint. The glenoid provides a surface for the articulation of the humerus, allowing for the wide range of motion in the shoulder. The stability and function of the shoulder joint are influenced by the shape and structure of the glenoid.';

api.ImageSearch mockImageSearch() {
  return api.ImageSearch(
    images: [
      "fixtures/glenoid_00.jpg",
      "fixtures/glenoid_01.jpg",
      "fixtures/glenoid_09.jpg",
      "fixtures/glenoid_11.jpg",
      "fixtures/glenoid_04.jpg",
      "fixtures/glenoid_14.jpg",
      "fixtures/glenoid_05.jpg",
      "fixtures/glenoid_12.jpg",
    ]
        .map((url) => api.SearchImage(
              contentUrl: url,
              contentSize: "",
              accentColor: "111111",
              encodingFormat: "jpg",
              height: 100,
              width: 100,
              hostPageUrl: "",
              thumbnailUrl: url,
              isLiked: false,
            ))
        .toList(),
    queryExpansions: [],
  );
}

class MyModel extends Model {
  bool _isConnected = false;
  WebSocketChannel? _channel;
  String _transcript = exampleTranscript;
  final List<api.ReferenceMaterial> _referenceMaterials = [];
  api.KeyTerms? _keyTerms = api.KeyTerms(
    entities: [
      //api.Entity(score: 1.0, text: "propranolol", type: "GENERIC_NAME"),
      //api.Entity(score: 1.0, text: "prednisone", type: "GENERIC_NAME"),
      //api.Entity(score: 1.0, text: "cancer", type: "DX_NAME"),
      //api.Entity(score: 1.0, text: "influenza", type: "DX_NAME"),
      //api.Entity(score: 1.0, text: "ativan", type: "BRAND_NAME"),
      //api.Entity(score: 1.0, text: "azithromycin", type: "GENERIC_NAME"),
      //api.Entity(score: 1.0, text: "sumatriptan", type: "GENERIC_NAME"),
      //api.Entity(score: 1.0, text: "parkinson's", type: "DX_NAME"),
      //api.Entity(score: 1.0, text: "eczema", type: "DX_NAME"),
      //api.Entity(score: 1.0, text: "parietal lobe", type: "SYSTEM_ORGAN_SITE"),
      //api.Entity(score: 1.0, text: "kidney", type: "SYSTEM_ORGAN_SITE"),
      //api.Entity(score: 1.0, text: "knee", type: "SYSTEM_ORGAN_SITE"),
      //api.Entity(score: 1.0, text: "left ventricle", type: "SYSTEM_ORGAN_SITE"),
      //api.Entity(score: 1.0, text: "broca's area", type: "SYSTEM_ORGAN_SITE"),
      api.Entity(score: 1.0, text: "glenoid", type: "SYSTEM_ORGAN_SITE"),
      api.Entity(score: 1.0, text: "humoral head", type: "SYSTEM_ORGAN_SITE"),
      api.Entity(score: 0.9, text: "osteophyte", type: "SYSTEM_ORGAN_SITE"),
      api.Entity(
          score: 1.0, text: "osteophyte crown", type: "SYSTEM_ORGAN_SITE"),
      api.Entity(score: 1.0, text: "humoral canal", type: "SYSTEM_ORGAN_SITE"),
      api.Entity(
          score: 1.0,
          text: "inferior glenoid capsule",
          type: "SYSTEM_ORGAN_SITE"),
      api.Entity(score: 1.0, text: "triceps tendon", type: "SYSTEM_ORGAN_SITE"),
      api.Entity(score: 0.9, text: "degenerative arthritis", type: "DX_NAME"),
      api.Entity(score: 0.8, text: "arthritis", type: "DX_NAME"),
      api.Entity(score: 0.7, text: "labral tear", type: "DX_NAME"),
      api.Entity(
          score: 1.0,
          text: "intra articular biceps tendon",
          type: "SYSTEM_ORGAN_SITE"),
      api.Entity(score: 0.65, text: "azithromycin", type: "GENERIC_NAME"),
    ],
  );
  api.Entity? _selectedEntity;
  api.ImageSearch? _searchImages;
  api.ImageSearch? _radiologySearchImages;
  api.ImageSearch? _histologySearchImages;

  final Map<String, api.DrugDetails?> _drugDetails = {};
  final Map<String, String> _definitions = {
    "glenoid": exampleDefinition,
  };
  final Set<String> _fetchingTerms = {};
  final Map<String, bool> _diseases = {};
  final Set<String> _fetchingDiseases = {};
  final Set<String> _fetchingDrugDetails = {};

  bool get isConnected => _isConnected;
  String get transcript => _transcript;
  List<api.ReferenceMaterial> get referenceMaterials => _referenceMaterials;
  api.KeyTerms? get keyTerms => _keyTerms;
  api.Entity? get selectedEntity => _selectedEntity;
  api.ImageSearch? get searchImages => _searchImages;
  api.ImageSearch? get radiologySearchImages => _radiologySearchImages;
  api.ImageSearch? get histologySearchImages => _histologySearchImages;

  void setDefinitionHelpful(bool helpful) {}

  Future<void> search(api.Entity entity) async {
    _searchImages = await api.search(
      entity.text,
      type: entity.type,
    );
    notifyListeners();
  }

  bool? hasDrugDetails(String query) {
    if (_drugDetails.containsKey(query)) return _drugDetails[query] != null;
    getDrugDetails(query);
    return null;
  }

  api.DrugDetails? getDrugDetails(String query) {
    if (_drugDetails.containsKey(query)) return _drugDetails[query];
    if (_fetchingDrugDetails.contains(query)) return null;
    _fetchingDrugDetails.add(query);
    api.getDrugDetails(query).then((details) {
      _drugDetails[query] = details;
      _fetchingDrugDetails.remove(query);
      notifyListeners();
    });
    return null;
  }

  bool? isDisease(String term) {
    final value = _diseases[term];
    if (value != null) return value;
    if (_fetchingDiseases.contains(term)) return null;
    _fetchingDiseases.add(term);
    api.isDisease(term).then((value) {
      _diseases[term] = value;
      _fetchingDiseases.remove(term);
      notifyListeners();
    });
    return null;
  }

  Future<void> searchRadiology(String text) async {
    _radiologySearchImages = await api.search(
      '$text radiology',
    );
    notifyListeners();
  }

  Future<void> searchHistology(String text) async {
    _histologySearchImages = await api.search(
      '$text histology',
    );
    notifyListeners();
  }

  set selectedEntity(api.Entity? entity) {
    if (_selectedEntity != entity) {
      _searchImages = null;
      _radiologySearchImages = null;
      _histologySearchImages = null;
    }
    _selectedEntity = entity;
    if (entity != null) {
      search(entity);
      switch (entity.type) {
        case 'DX_NAME':
          searchHistology(entity.text);
          break;
        case 'SYSTEM_ORGAN_SITE':
          searchRadiology(entity.text);
          break;
        default:
          break;
      }
    }
    notifyListeners();
  }

  MyModel() {
    connectWebSock();
    _selectedEntity =
        _keyTerms?.entities.firstWhere((e) => e.text == "glenoid");
    _searchImages = mockImageSearch();
  }

  void onConnect() {
    _isConnected = true;
    notifyListeners();
  }

  void onDisconnect() {
    _isConnected = false;
    notifyListeners();
  }

  bool termIsFetching(String term) => _fetchingTerms.contains(term);

  String? define(String term) {
    final value = _definitions[term];
    if (value != null) return value;
    if (_fetchingTerms.contains(term)) return null;
    _fetchingTerms.add(term);
    api.define(term).then((value) {
      _definitions[term] = value;
      _fetchingTerms.remove(term);
      notifyListeners();
    }).catchError((err) {
      print('error fetching definition for $term: $err');
    });
    return null;
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
