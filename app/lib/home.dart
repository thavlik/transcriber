import 'dart:math';
import 'dart:typed_data';

import 'package:demo/pharmacy.dart';
import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:scoped_model/scoped_model.dart';

import 'model.dart';
import 'api.dart';

final kTransparentImageBytes = Uint8List.fromList(<int>[
  0x89,
  0x50,
  0x4E,
  0x47,
  0x0D,
  0x0A,
  0x1A,
  0x0A,
  0x00,
  0x00,
  0x00,
  0x0D,
  0x49,
  0x48,
  0x44,
  0x52,
  0x00,
  0x00,
  0x00,
  0x01,
  0x00,
  0x00,
  0x00,
  0x01,
  0x08,
  0x06,
  0x00,
  0x00,
  0x00,
  0x1F,
  0x15,
  0xC4,
  0x89,
  0x00,
  0x00,
  0x00,
  0x0A,
  0x49,
  0x44,
  0x41,
  0x54,
  0x78,
  0x9C,
  0x63,
  0x00,
  0x01,
  0x00,
  0x00,
  0x05,
  0x00,
  0x01,
  0x0D,
  0x0A,
  0x2D,
  0xB4,
  0x87,
  0x9C,
  0x00,
  0x00,
  0x00,
  0x00,
  0x49,
  0x45,
  0x4E,
  0x44,
  0xAE,
]);
final kTransparentImage = Image.memory(kTransparentImageBytes);

extension HexColor on Color {
  /// String is in the format "aabbcc" or "ffaabbcc" with an optional leading "#".
  static Color fromHex(String hexString) {
    final buffer = StringBuffer();
    if (hexString.length == 6 || hexString.length == 7) buffer.write('ff');
    buffer.write(hexString.replaceFirst('#', ''));
    return Color(int.parse(buffer.toString(), radix: 16));
  }

  /// Prefixes a hash sign if [leadingHashSign] is set to `true` (default is `true`).
  String toHex({bool leadingHashSign = true}) => '${leadingHashSign ? '#' : ''}'
      '${alpha.toRadixString(16).padLeft(2, '0')}'
      '${red.toRadixString(16).padLeft(2, '0')}'
      '${green.toRadixString(16).padLeft(2, '0')}'
      '${blue.toRadixString(16).padLeft(2, '0')}';
}

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

class _HomePageState extends State<HomePage> with TickerProviderStateMixin {
  String? _viewImage;

  void onImageTap(
    BuildContext context,
    String imageUrl,
  ) =>
      setState(() => _viewImage = imageUrl);

  void _onEntitySelected(BuildContext context, Entity entity) {
    final model = ScopedModel.of<MyModel>(context);
    if (model.selectedEntity == entity) {
      // deselect entity
      model.selectedEntity = null;
      return;
    }
    model.selectedEntity = entity;
  }

  Widget _buildEntity(BuildContext context, Entity entity) => InkWell(
        onTap: () => _onEntitySelected(context, entity),
        child: Container(
          decoration: BoxDecoration(
            color: ScopedModel.of<MyModel>(context).selectedEntity?.text ==
                    entity.text
                ? Theme.of(context).highlightColor
                : null,
            border: Border(
              bottom: BorderSide(
                color: Theme.of(context).dividerColor,
                width: 1,
              ),
            ),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Flexible(
                child: Padding(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 8.0,
                    vertical: 10.0,
                  ),
                  child: Text(
                    entity.text,
                    overflow: TextOverflow.ellipsis,
                    style: Theme.of(context).textTheme.labelMedium,
                  ),
                ),
              ),
              Flexible(
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Text(
                    entity.type.replaceAll('_', ' '),
                    overflow: TextOverflow.ellipsis,
                    style: Theme.of(context).textTheme.labelSmall!.copyWith(
                        color: HSLColor.fromAHSL(
                                1.0, entity.score * 180.0, 0.5, 0.5)
                            .toColor()),
                  ),
                ),
              ),
              //Padding(
              //  padding: const EdgeInsets.all(8.0),
              //  child: Text(
              //    ' (${entity.score})',
              //    style: Theme.of(context).textTheme.bodyLarge,
              //  ),
              //),
            ],
          ),
        ),
      );

  Widget _buildKeyTerms(
    BuildContext context,
    KeyTerms keyTerms,
    Entity? selectedEntity,
  ) {
    List<Entity> entities = [...keyTerms.entities];
    if (selectedEntity != null) {
      if (entities.indexWhere((e) => e.text == selectedEntity.text) == -1) {
        entities.add(selectedEntity);
      }
      entities.sort((a, b) => a.text.compareTo(b.text));
    }
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: entities.map((e) => _buildEntity(context, e)).toList(),
    );
  }

  Widget _columnHeader(BuildContext context, String name) => Flexible(
        flex: 1,
        child: Text(name, style: Theme.of(context).textTheme.headlineSmall),
      );

  Widget _buildSearchImage(
    BuildContext context,
    SearchImage img,
    MyModel model,
  ) =>
      Padding(
        padding: const EdgeInsets.all(8.0),
        child: Container(
          constraints: const BoxConstraints(
            maxWidth: 100,
            maxHeight: 100,
          ),
          child: Center(
            child: InkWell(
              onTap: () => onImageTap(context, img.contentUrl),
              child: Container(
                  decoration: BoxDecoration(
                    color: HexColor.fromHex(img.accentColor),
                  ),
                  child: Stack(
                    children: [
                      AspectRatio(
                        aspectRatio: img.width / img.height,
                        child: FadeInImage(
                          placeholder: MemoryImage(kTransparentImageBytes),
                          image: NetworkImage(img.thumbnailUrl),
                          fit: BoxFit.cover,
                          imageErrorBuilder: (context, error, stackTrace) {
                            return kTransparentImage;
                            //return Image.asset('assets/images/error.jpg',
                            //    fit: BoxFit.fitWidth);
                          },
                        ),
                      ),
                      Positioned(
                        top: 3,
                        right: 3,
                        child: Opacity(
                          opacity: img.isLiked == true ? 0.8 : 0.7,
                          child: GestureDetector(
                            child: Container(
                              decoration: BoxDecoration(
                                boxShadow: BoxDecoration(
                                  boxShadow: [
                                    BoxShadow(
                                      color: Colors.black.withAlpha(80),
                                      blurRadius: 8.0,
                                      spreadRadius: 0.0,
                                      offset: Offset.zero,
                                    ),
                                  ],
                                ).boxShadow,
                              ),
                              child: Icon(
                                img.isLiked == true
                                    ? Icons.favorite
                                    : Icons.favorite_border,
                              ),
                            ),
                            onTap: () =>
                                model.likeImage(img, !(img.isLiked ?? false)),
                          ),
                        ),
                      ),
                    ],
                  )),
            ),
          ),
        ),
      );

  Widget _pharmacologyTab(BuildContext context, MyModel model) {
    if (model.selectedEntity == null) return Container();
    final has = model.hasDrugDetails(model.selectedEntity!.text);
    if (has == null) {
      return const Center(
        child: CircularProgressIndicator(),
      );
    } else if (has == false) {
      return Container();
    }
    final details = model.getDrugDetails(model.selectedEntity!.text)!;
    return PharmacyView(
      details,
      onImageTap: (url) => onImageTap(context, url),
    );
  }

  Widget _buildRefTab(BuildContext context, MyModel model) {
    final selectedEntity = model.selectedEntity;
    final showRadiologyTab = selectedEntity?.type == 'SYSTEM_ORGAN_SITE';
    final showHistologyTab = selectedEntity != null &&
        selectedEntity.type == 'DX_NAME' &&
        model.isDisease(selectedEntity.text) == true;
    final showPharmacologyTab = selectedEntity?.type == "GENERIC_NAME" ||
        selectedEntity?.type == "BRAND_NAME";
    int numTabs = 1;
    if (showRadiologyTab) numTabs++;
    if (showHistologyTab) numTabs++;
    if (showPharmacologyTab) numTabs++;
    return DefaultTabController(
      length: numTabs,
      child: Builder(builder: (context) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            TabBar(
              tabs: [
                const Tab(text: 'Overview'),
                if (showRadiologyTab) const Tab(text: 'Radiology'),
                if (showHistologyTab) const Tab(text: 'Histology'),
                if (showPharmacologyTab) const Tab(text: 'Pharma'),
              ],
            ),
            Expanded(
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: TabBarView(
                  children: [
                    model.selectedEntity == null
                        ? Container()
                        : SingleChildScrollView(
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.start,
                              crossAxisAlignment: CrossAxisAlignment.stretch,
                              children: [
                                Text(
                                  'Definition of "${model.selectedEntity!.text}":',
                                  style: Theme.of(context).textTheme.labelLarge,
                                ),
                                Padding(
                                  padding: const EdgeInsets.all(8.0),
                                  child: Container(
                                    decoration: BoxDecoration(
                                      border: Border.all(
                                        color: Theme.of(context).dividerColor,
                                      ),
                                      borderRadius: BorderRadius.circular(8),
                                    ),
                                    child: Padding(
                                        padding: const EdgeInsets.all(4.0),
                                        child: model.termIsFetching(
                                                model.selectedEntity!.text)
                                            ? const Padding(
                                                padding: EdgeInsets.all(16.0),
                                                child: Center(
                                                    child:
                                                        CircularProgressIndicator()),
                                              )
                                            : Padding(
                                                padding:
                                                    const EdgeInsets.all(8.0),
                                                child: Text(
                                                  model.define(model
                                                          .selectedEntity!
                                                          .text) ??
                                                      '',
                                                ),
                                              )),
                                  ),
                                ),
                                if (model.define(model.selectedEntity!.text) !=
                                    null)
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                      horizontal: 16,
                                    ),
                                    child: Column(
                                      crossAxisAlignment:
                                          CrossAxisAlignment.end,
                                      children: [
                                        Text(
                                          'Was this definition helpful?',
                                          style: Theme.of(context)
                                              .textTheme
                                              .bodySmall,
                                        ),
                                        Row(
                                          mainAxisAlignment:
                                              MainAxisAlignment.end,
                                          children: [
                                            TextButton(
                                              onPressed: () => model
                                                  .setDefinitionHelpful(true),
                                              child: const Text('Yes'),
                                            ),
                                            TextButton(
                                              onPressed: () => model
                                                  .setDefinitionHelpful(false),
                                              child: const Text('No'),
                                            ),
                                          ],
                                        ),
                                      ],
                                    ),
                                  ),
                                const SizedBox(height: 16),
                                Column(
                                  crossAxisAlignment:
                                      CrossAxisAlignment.stretch,
                                  children: [
                                    Text(
                                      'Images of "${model.selectedEntity!.text}":',
                                      style: Theme.of(context)
                                          .textTheme
                                          .labelLarge,
                                    ),
                                    model.searchImages == null
                                        ? Container()
                                        : Wrap(
                                            children: model
                                                .searchImages!.queryExpansions
                                                .sublist(
                                                    0,
                                                    model.searchImages!
                                                        .queryExpansions.length)
                                                .map((e) => QueryExpansion(
                                                      model,
                                                      model.selectedEntity!,
                                                      e,
                                                    ))
                                                .toList(),
                                          ),
                                  ],
                                ),
                                Padding(
                                  padding: const EdgeInsets.all(8.0),
                                  child: Container(
                                    decoration: BoxDecoration(
                                      border: Border.all(
                                        color: Theme.of(context).dividerColor,
                                      ),
                                      borderRadius: BorderRadius.circular(8),
                                    ),
                                    child: model.searchImages == null
                                        ? const Padding(
                                            padding: EdgeInsets.all(16.0),
                                            child: Center(
                                              child:
                                                  CircularProgressIndicator(),
                                            ),
                                          )
                                        : Column(
                                            crossAxisAlignment:
                                                CrossAxisAlignment.stretch,
                                            children: [
                                              Wrap(
                                                crossAxisAlignment:
                                                    WrapCrossAlignment.start,
                                                alignment:
                                                    WrapAlignment.spaceAround,
                                                direction: Axis.horizontal,
                                                children: model
                                                    .searchImages!.images
                                                    .map((e) =>
                                                        _buildSearchImage(
                                                            context, e, model))
                                                    .toList(),
                                              ),
                                              TextButton(
                                                  onPressed: () {},
                                                  child: const Padding(
                                                    padding:
                                                        EdgeInsets.all(8.0),
                                                    child: Text('Load more...'),
                                                  )),
                                            ],
                                          ),
                                  ),
                                ),
                                Padding(
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 16,
                                  ),
                                  child: Column(
                                    crossAxisAlignment: CrossAxisAlignment.end,
                                    children: [
                                      Text(
                                        'Are these images helpful?',
                                        style: Theme.of(context)
                                            .textTheme
                                            .bodySmall,
                                      ),
                                      Row(
                                        mainAxisAlignment:
                                            MainAxisAlignment.end,
                                        children: [
                                          TextButton(
                                            onPressed: () => model
                                                .setDefinitionHelpful(true),
                                            child: const Text('Yes'),
                                          ),
                                          TextButton(
                                            onPressed: () => model
                                                .setDefinitionHelpful(false),
                                            child: const Text('No'),
                                          ),
                                        ],
                                      ),
                                    ],
                                  ),
                                ),
                              ],
                            ),
                          ),
                    if (showRadiologyTab)
                      model.selectedEntity == null
                          ? Container()
                          : model.radiologySearchImages == null
                              ? const Center(
                                  child: CircularProgressIndicator(),
                                )
                              : SingleChildScrollView(
                                  child: Column(
                                    children: [
                                      Wrap(
                                        crossAxisAlignment:
                                            WrapCrossAlignment.start,
                                        alignment: WrapAlignment.spaceAround,
                                        direction: Axis.horizontal,
                                        children: model
                                            .radiologySearchImages!.images
                                            .map((e) => _buildSearchImage(
                                                context, e, model))
                                            .toList(),
                                      ),
                                      TextButton(
                                          onPressed: () {},
                                          child: const Padding(
                                            padding: EdgeInsets.all(8.0),
                                            child: Text('Load more...'),
                                          ))
                                    ],
                                  ),
                                ),
                    if (showHistologyTab)
                      model.selectedEntity == null
                          ? Container()
                          : model.histologySearchImages == null
                              ? const Center(
                                  child: CircularProgressIndicator(),
                                )
                              : SingleChildScrollView(
                                  child: Column(
                                    children: [
                                      Wrap(
                                        crossAxisAlignment:
                                            WrapCrossAlignment.start,
                                        alignment: WrapAlignment.spaceAround,
                                        direction: Axis.horizontal,
                                        children: model
                                            .histologySearchImages!.images
                                            .map((e) => _buildSearchImage(
                                                context, e, model))
                                            .toList(),
                                      ),
                                      TextButton(
                                          onPressed: () {},
                                          child: const Padding(
                                            padding: EdgeInsets.all(8.0),
                                            child: Text('Load more...'),
                                          )),
                                      Padding(
                                        padding: const EdgeInsets.symmetric(
                                          horizontal: 16,
                                        ),
                                        child: Column(
                                          crossAxisAlignment:
                                              CrossAxisAlignment.end,
                                          children: [
                                            Text(
                                              'Are these images helpful?',
                                              style: Theme.of(context)
                                                  .textTheme
                                                  .bodySmall,
                                            ),
                                            Row(
                                              mainAxisAlignment:
                                                  MainAxisAlignment.end,
                                              children: [
                                                TextButton(
                                                  onPressed: () => model
                                                      .setDefinitionHelpful(
                                                          true),
                                                  child: const Text('Yes'),
                                                ),
                                                TextButton(
                                                  onPressed: () => model
                                                      .setDefinitionHelpful(
                                                          false),
                                                  child: const Text('No'),
                                                ),
                                              ],
                                            ),
                                          ],
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                    if (showPharmacologyTab) _pharmacologyTab(context, model),
                  ],
                ),
              ),
            ),
          ],
        );
      }),
    );
  }

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
          body: Stack(
            children: [
              Column(
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
                  Expanded(
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Flexible(
                          flex: 1,
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.start,
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Flexible(
                                child: SingleChildScrollView(
                                  child: Padding(
                                    padding: const EdgeInsets.all(8.0),
                                    child: Text(model.transcript,
                                        style: Theme.of(context)
                                            .textTheme
                                            .bodyLarge!
                                            .copyWith(
                                              color: Colors.deepOrange
                                                  .withOpacity(0.75),
                                            )),
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                        Flexible(
                          flex: 1,
                          child: model.keyTerms == null
                              ? Container()
                              : SingleChildScrollView(
                                  child: _buildKeyTerms(
                                    context,
                                    model.keyTerms!,
                                    model.selectedEntity,
                                  ),
                                ),
                        ),
                        Flexible(
                          flex: 1,
                          child: _buildRefTab(context, model),
                        ),
                      ],
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
                        child: _viewImage!.contains('.svg')
                            ? SvgPicture.network(_viewImage!)
                            : FadeInImage(
                                placeholder:
                                    MemoryImage(kTransparentImageBytes),
                                image: NetworkImage(_viewImage!),
                                fit: BoxFit.cover,
                                imageErrorBuilder:
                                    (context, error, stackTrace) {
                                  return kTransparentImage;
                                  //return Image.asset('assets/images/error.jpg',
                                  //    fit: BoxFit.fitWidth);
                                },
                              ),
                      ),
                    ),
                  ),
                ),
            ],
          ),
        );
      },
    );
  }
}

class QueryExpansion extends StatelessWidget {
  final MyModel model;
  final String query;
  final Entity parent;
  final Color? color;

  const QueryExpansion(
    this.model,
    this.parent,
    this.query, {
    super.key,
    this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(
        horizontal: 4,
        vertical: 2,
      ),
      child: Container(
        decoration: BoxDecoration(
          color: color ?? Colors.purple.shade900,
          border: Border.all(
            color: Theme.of(context).dividerColor,
            width: 1,
          ),
          borderRadius: BorderRadius.circular(16),
        ),
        child: InkWell(
          onTap: () => model.selectedEntity = Entity(
            score: 1.0,
            text: query,
            type: parent.type,
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(
              vertical: 4.0,
              horizontal: 8.0,
            ),
            child: Text(
              query,
              style: Theme.of(context).textTheme.bodySmall!.copyWith(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                  ),
            ),
          ),
        ),
      ),
    );
  }
}
