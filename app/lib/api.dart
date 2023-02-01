import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';

const apiHost = 'ts.beebs.dev';

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
  final url = Uri.https(apiHost, '/api/search', {
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
