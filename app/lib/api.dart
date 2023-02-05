import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';

const apiHost = 'ts.beebs.dev';

class DrugStructure {
  final String image;
  final String pdb;

  DrugStructure({
    required this.image,
    required this.pdb,
  });

  factory DrugStructure.fromJson(Map<String, dynamic> json) {
    return DrugStructure(
      image: json['image'] as String,
      pdb: json['pDB'] as String,
    );
  }
}

class DrugWeight {
  final String average;
  final String monoisotopic;

  DrugWeight({
    required this.average,
    required this.monoisotopic,
  });

  factory DrugWeight.fromJson(Map<String, dynamic> json) {
    return DrugWeight(
      average: json['average'] as String,
      monoisotopic: json['monoisotopic'] as String,
    );
  }
}

class DrugPharmacology {
  final String indication;
  final List<String> associatedConditions;
  final String pharmacodynamics;
  final String mechanismOfAction;
  final String absorption;
  final String volumeOfDistribution;
  final String proteinBinding;
  final DrugMetabolism metabolism;
  final String routeOfElimination;
  final String halfLife;
  final String clearance;
  final String toxicity;

  DrugPharmacology({
    required this.indication,
    required this.associatedConditions,
    required this.pharmacodynamics,
    required this.mechanismOfAction,
    required this.absorption,
    required this.volumeOfDistribution,
    required this.proteinBinding,
    required this.metabolism,
    required this.routeOfElimination,
    required this.halfLife,
    required this.clearance,
    required this.toxicity,
  });

  factory DrugPharmacology.fromJson(Map<String, dynamic> json) {
    return DrugPharmacology(
      indication: json['indication'] as String,
      associatedConditions:
          (json['associatedConditions'] as List<dynamic>).cast<String>(),
      pharmacodynamics: json['pharmacodynamics'] as String,
      mechanismOfAction: json['mechanismOfAction'] as String,
      absorption: json['absorption'] as String,
      volumeOfDistribution: json['volumeOfDistribution'] as String,
      proteinBinding: json['proteinBinding'] as String,
      metabolism:
          DrugMetabolism.fromJson(json['metabolism'] as Map<String, dynamic>),
      routeOfElimination: json['routeOfElimination'] as String,
      halfLife: json['halfLife'] as String,
      clearance: json['clearance'] as String,
      toxicity: json['toxicity'] as String,
    );
  }
}

class DrugReference {
  final int index;
  final String title;
  final String link;

  DrugReference({
    required this.index,
    required this.title,
    required this.link,
  });

  factory DrugReference.fromJson(Map<String, dynamic> json) {
    return DrugReference(
      index: json['index'] as int,
      title: json['title'] as String,
      link: json['link'] as String,
    );
  }
}

class DrugReferences {
  final List<DrugReference> general;

  DrugReferences({required this.general});

  factory DrugReferences.fromJson(Map<String, dynamic> json) {
    return DrugReferences(
      general: (json['general'] as List<dynamic>)
          .map((e) => DrugReference.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }
}

class DrugMetabolism {
  final String description;

  DrugMetabolism({required this.description});

  factory DrugMetabolism.fromJson(Map<String, dynamic> json) {
    return DrugMetabolism(
      description: json['description'] as String,
    );
  }
}

class DrugDetails {
  final String summary;
  final List<String> brandNames;
  final String genericName;
  final String drugBankAccessionNumber;
  final String background;
  final String type;
  final List<String> groups;
  final DrugStructure structure;
  final DrugWeight weight;
  final String chemicalFormula;
  final List<String> synonyms;
  final List<String> externalIDs;
  final DrugPharmacology pharmacology;
  final DrugReferences references;

  DrugDetails({
    required this.summary,
    required this.brandNames,
    required this.genericName,
    required this.type,
    required this.drugBankAccessionNumber,
    required this.background,
    required this.groups,
    required this.structure,
    required this.weight,
    required this.chemicalFormula,
    required this.synonyms,
    required this.externalIDs,
    required this.pharmacology,
    required this.references,
  });

  factory DrugDetails.fromJson(Map<String, dynamic> json) {
    return DrugDetails(
      summary: json['summary'] as String,
      type: json['type'] as String,
      brandNames: (json['brandNames'] as List<dynamic>).cast<String>(),
      genericName: json['genericName'] as String,
      drugBankAccessionNumber: json['drugBankAccessionNumber'] as String,
      background: json['background'] as String,
      groups: (json['groups'] as List<dynamic>).cast<String>(),
      structure:
          DrugStructure.fromJson(json['structure'] as Map<String, dynamic>),
      weight: DrugWeight.fromJson(json['weight'] as Map<String, dynamic>),
      chemicalFormula: json['chemicalFormula'] as String,
      synonyms: (json['synonyms'] as List<dynamic>).cast<String>(),
      externalIDs: (json['externalIDs'] as List<dynamic>).cast<String>(),
      pharmacology: DrugPharmacology.fromJson(
          json['pharmacology'] as Map<String, dynamic>),
      references:
          DrugReferences.fromJson(json['references'] as Map<String, dynamic>),
    );
  }
}

class Entity {
  String text;
  String type;
  double score;

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
  final Map<Entity, DateTime> _lastUsed = {};
  final Duration pruneAfter;

  KeyTerms({
    required this.entities,
    this.pruneAfter = const Duration(seconds: 3),
  });

  void sort() => entities.sort((a, b) => a.text.compareTo(b.text));

  void integrate(List<Entity> other) {
    for (var entity in other) {
      final lower = entity.text.toLowerCase();
      bool found = false;
      for (var existing in entities) {
        if (existing.text.toLowerCase() == lower &&
            existing.type == entity.type) {
          // update the score
          existing.score = entity.score;
          used(existing);
          found = true;
          break;
        }
      }
      if (!found) {
        entities.add(entity);
        used(entity);
      }
    }
    sort();
  }

  void prune() {
    final List<Entity> remove = [];
    for (var entity in entities) {
      final used = _lastUsed[entity];
      if (used == null) {
        continue;
      }
      if (DateTime.now().difference(used) > pruneAfter) {
        remove.add(entity);
        _lastUsed.remove(entity);
      }
    }
    for (var entity in remove) {
      entities.remove(entity);
    }
  }

  void used(Entity e) => _lastUsed[e] = DateTime.now();

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

Future<DrugDetails?> getDrugDetails(String query) async {
  final url = Uri.https(apiHost, '/drug', {
    'q': query,
  });
  final response = await http.get(url);
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as Map<String, dynamic>;
  return DrugDetails.fromJson(decodedResponse);
}

Future<String> define(String query) async {
  final url = Uri.https(apiHost, '/define', {
    'q': query,
  });
  final response = await http.get(url);
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as Map<String, dynamic>;
  return decodedResponse["text"] as String;
}

Future<bool> isDisease(String query) async {
  final url = Uri.https(apiHost, '/disease', {
    'q': query,
  });
  final response = await http.get(url);
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as Map<String, dynamic>;
  return decodedResponse["isDisease"] as bool;
}

Future<List<SearchImage>> search(String query) async {
  final url = Uri.https(apiHost, '/img/search', {
    'q': query,
  });
  final response = await http.get(url);
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as List? ?? [];
  return decodedResponse.map((e) => SearchImage.fromJson(e)).toList();
}

Future<void> likeImage(SearchImage img, bool isLiked) async {
  final url = Uri.https(apiHost, '/like', {
    'i': '',
  });
  final response = await http.post(url,
      body: json.encode({
        'hash': '',
        'isLiked': isLiked,
      }));
  checkHttpStatus(response);
}
