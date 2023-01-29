import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import 'model.dart';

class ReferenceMaterialWidget extends StatelessWidget {
  const ReferenceMaterialWidget(
    this.model, {
    super.key,
    required this.onImageTap,
  });

  final ReferenceMaterial model;
  final Function(ReferenceMaterial, String) onImageTap;

  Widget _buildTerm(BuildContext context, String term) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.all(8.0),
          child: Text(
            term,
            style: Theme.of(context).textTheme.headlineSmall,
          ),
        ),
        Padding(
          padding: const EdgeInsets.all(8.0),
          child: Text(
            '(${model.terms[0]})',
            style: Theme.of(context).textTheme.bodyMedium,
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
          border: Border(
        bottom: BorderSide(
          color: Theme.of(context).dividerColor,
          width: 1,
        ),
      )),
      child: Row(
        children: [
          Expanded(
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                _buildTerm(context, model.matched),
                Row(
                  children: model.images
                      .map((e) => Padding(
                            padding: const EdgeInsets.all(8.0),
                            child: InkWell(
                              onTap: () => onImageTap(model, e),
                              child: Container(
                                width: 100,
                                height: 100,
                                decoration: BoxDecoration(
                                  image: DecorationImage(
                                    image: NetworkImage(e),
                                    fit: BoxFit.cover,
                                  ),
                                ),
                              ),
                            ),
                          ))
                      .toList(),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  String? _viewImage;

  void onImageTap(
    BuildContext context,
    ReferenceMaterial ref,
    String imageUrl,
  ) =>
      setState(() => _viewImage = imageUrl);

  Widget _buildEntity(BuildContext context, Entity entity) => Row(
        children: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(
              entity.text,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(
              ' (${entity.type})',
              style: Theme.of(context).textTheme.bodyLarge,
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(
              ' (${entity.score})',
              style: Theme.of(context).textTheme.bodyLarge,
            ),
          ),
        ],
      );

  Widget _buildKeyTerms(BuildContext context, KeyTerms keyTerms) => Container(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children:
              keyTerms.entities.map((e) => _buildEntity(context, e)).toList(),
        ),
      );

  Widget _columnHeader(BuildContext context, String name) => Flexible(
        flex: 1,
        child: Text(name, style: Theme.of(context).textTheme.headlineSmall),
      );

  @override
  Widget build(BuildContext context) {
    return ScopedModelDescendant<MyModel>(
      builder: (context, child, model) {
        return Scaffold(
          appBar: AppBar(
            title: Text(
              model.isConnected ? 'Connected' : 'Not Connected',
            ),
          ),
          body: Container(
            constraints: const BoxConstraints.expand(),
            child: Stack(
              children: [
                SingleChildScrollView(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      Container(
                        decoration: BoxDecoration(
                          border: Border(
                            bottom: BorderSide(
                              color: Theme.of(context).dividerColor,
                              width: 1,
                            ),
                          ),
                        ),
                        child: Padding(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 16.0,
                            vertical: 8.0,
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              _columnHeader(context, 'Transcript'),
                              _columnHeader(context, 'Key Terms'),
                              _columnHeader(context, 'Reference Materials'),
                            ],
                          ),
                        ),
                      ),
                      Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Flexible(
                              flex: 1,
                              child: Row(
                                mainAxisAlignment: MainAxisAlignment.start,
                                children: [
                                  Flexible(
                                    child: Padding(
                                      padding: const EdgeInsets.all(8.0),
                                      child: Text(model.transcript ?? '',
                                          style: Theme.of(context)
                                              .textTheme
                                              .bodyLarge!
                                              .copyWith(
                                                color: Colors.deepOrange
                                                    .withOpacity(0.75),
                                              )),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            Flexible(
                              flex: 1,
                              child: model.keyTerms != null
                                  ? _buildKeyTerms(context, model.keyTerms!)
                                  : Container(),
                            ),
                            Flexible(
                              flex: 1,
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: model.referenceMaterials.reversed
                                    .map((ref) => ReferenceMaterialWidget(
                                          ref,
                                          onImageTap: (ref, img) =>
                                              onImageTap(context, ref, img),
                                        ))
                                    .toList(),
                              ),
                            ),
                          ]),
                    ],
                  ),
                ),
                if (_viewImage != null)
                  Positioned.fill(
                    child: GestureDetector(
                      onTap: () => setState(() => _viewImage = null),
                      child: Container(
                        color: Colors.black.withOpacity(0.5),
                        child: Center(
                          child: Container(
                            width: 800,
                            height: 600,
                            decoration: BoxDecoration(
                              image: DecorationImage(
                                image: NetworkImage(_viewImage!),
                                fit: BoxFit.contain,
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
              ],
            ),
          ),
        );
      },
    );
  }
}
