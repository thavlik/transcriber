import 'package:demo/api.dart';
import 'package:demo/structure_view.dart';
import 'package:expandable/expandable.dart';
import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:url_launcher/url_launcher.dart';

import 'home.dart';

class PharmacyView extends StatefulWidget {
  final DrugDetails drug;

  final Function(String imageUrl) onImageTap;

  const PharmacyView(
    this.drug, {
    super.key,
    required this.onImageTap,
  });

  @override
  State<PharmacyView> createState() => _PharmacyViewState();
}

class _PharmacyViewState extends State<PharmacyView> {
  final ExpandableController _general = ExpandableController();
  final ExpandableController _pharmacology = ExpandableController();
  final ExpandableController _references = ExpandableController();

  DrugDetails get drug => widget.drug;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      setState(() => _general.expanded = true);
    });
  }

  Widget _sectionHeader(String title) => Text(
        title,
        style: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.bold,
        ),
      );

  Widget _section(String title, Widget child) => Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _sectionHeader(title),
            child,
          ],
        ),
      );

  Widget _buildWeight() => _section(
      'Weight',
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Average: ${drug.weight.average}'),
          Text('Monoisotopic: ${drug.weight.monoisotopic}'),
        ],
      ));

  Widget _buildExpandableHeader(String title, bool collapsed) => Container(
        color: Theme.of(context).cardColor,
        child: ExpandableButton(
          // <-- Collapses when tapped on
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Text(
                  title,
                  style: Theme.of(context).textTheme.bodyMedium,
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Icon(collapsed ? Icons.add : Icons.remove),
              ),
            ],
          ),
        ),
      );

  Widget _buildReferences(BuildContext context) => _section(
      'General',
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          for (final reference in drug.references.general)
            InkWell(
              onTap: () {
                if (reference.link == null) {
                  return;
                }
                final url = Uri.parse(reference.link!);
                launchUrl(url);
              },
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    SizedBox(
                      width: 32,
                      child: Text(
                        '[${reference.index}]',
                        style: Theme.of(context).textTheme.bodyMedium,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Flexible(
                      child: Text(
                        '${reference.title} ${reference.link}',
                        style: Theme.of(context).textTheme.bodySmall,
                      ),
                    ),
                  ],
                ),
              ),
            ),
        ],
      ));

  Widget _buildStructureView(BuildContext context) {
    final structureUrl =
        'https://ts.beebs.dev/drug/${drug.drugBankAccessionNumber}/structure.svg';
    return Container(
      constraints: const BoxConstraints(
        maxHeight: 300,
        maxWidth: 300,
      ),
      child: DefaultTabController(
          length: 2,
          child: Column(
            children: [
              Expanded(
                child: Container(
                  decoration: BoxDecoration(
                    border: Border.all(color: Colors.white),
                  ),
                  child: TabBarView(
                    children: [
                      GestureDetector(
                        onTap: () => widget.onImageTap(structureUrl),
                        child: Container(
                          color: Colors.white,
                          width: 100,
                          height: 100,
                          child: SvgPicture.network(structureUrl),
                        ),
                      ),
                      StructureView(drug),
                    ],
                  ),
                ),
              ),
              const TabBar(
                tabs: [
                  Tab(text: 'Structure'),
                  Tab(text: '3D'),
                ],
              ),
            ],
          )),
    );
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ExpandablePanel(
            controller: _general,
            collapsed: _buildExpandableHeader('General', true),
            expanded: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildExpandableHeader('General', false),
                _section('Summary', Text(drug.summary)),
                _buildStructureView(context),
                _section('Brand Names', Text(drug.brandNames.join(', '))),
                _section('Generic Name', Text(drug.genericName)),
                _section('Background', Text(drug.background)),
                _section('Groups', Text(drug.groups.join(', '))),
                _section('Type', Text(drug.type)),
                _section('Chemical Formula', Text(drug.chemicalFormula)),
                _buildWeight(),
              ],
            ),
          ),
          const SizedBox(height: 8),
          ExpandablePanel(
            controller: _pharmacology,
            collapsed: _buildExpandableHeader('Pharmacology', true),
            expanded: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildExpandableHeader('Pharmacology', false),
                _section('Indication', Text(drug.pharmacology.indication)),
                _section('Associated conditions',
                    Text(drug.pharmacology.associatedConditions.join(', '))),
                _section('Pharmacodynamics',
                    Text(drug.pharmacology.pharmacodynamics)),
                _section('Mechanism of action',
                    Text(drug.pharmacology.mechanismOfAction.split('\n')[0])),
                _section('Absorption', Text(drug.pharmacology.absorption)),
                _section('Volume of distribution',
                    Text(drug.pharmacology.volumeOfDistribution)),
                _section('Metabolism',
                    Text(drug.pharmacology.metabolism.description)),
                _section('Route of elimination',
                    Text(drug.pharmacology.routeOfElimination)),
                _section('Half-life', Text(drug.pharmacology.halfLife)),
                _section('Clearance', Text(drug.pharmacology.clearance)),
                _section('Toxicity', Text(drug.pharmacology.toxicity)),
              ],
            ),
          ),
          const SizedBox(height: 8),
          ExpandablePanel(
            controller: _references,
            collapsed: _buildExpandableHeader('References', true),
            expanded: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildExpandableHeader('References', false),
                _buildReferences(context),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
