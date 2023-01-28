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

  Widget _buildTerm(String term) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Text(term),
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
                _buildTerm(model.matched),
                //Column(
                //  crossAxisAlignment: CrossAxisAlignment.start,
                //  children: [
                //    ...model.terms.map((e) => _buildTerm(e)).toList(),
                //  ],
                //),
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

  @override
  Widget build(BuildContext context) {
    return ScopedModelDescendant<MyModel>(
      builder: (context, child, model) {
        return Scaffold(
          appBar: AppBar(
            title: Text(model.isConnected ? 'Connected' : 'Not Connected'),
          ),
          body: Container(
            constraints: const BoxConstraints.expand(),
            child: Stack(
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    SingleChildScrollView(
                      child: SizedBox(
                        width: 1000,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: model.referenceMaterials
                              .map((ref) => ReferenceMaterialWidget(
                                    ref,
                                    onImageTap: (ref, img) =>
                                        onImageTap(context, ref, img),
                                  ))
                              .toList(),
                        ),
                      ),
                    ),
                  ],
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
