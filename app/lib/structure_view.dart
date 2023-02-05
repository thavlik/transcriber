import 'package:demo/parse_bstl.dart';
import 'package:flutter/material.dart';
import 'package:simple_3d/simple_3d.dart';
import 'package:util_simple_3d/util_simple_3d.dart';
import 'package:simple_3d_renderer/simple_3d_renderer.dart';
import 'api.dart' as api;

class StructureView extends StatefulWidget {
  final api.DrugDetails drugDetails;

  const StructureView(this.drugDetails, {Key? key}) : super(key: key);

  @override
  State<StatefulWidget> createState() => _StructureViewState();
}

class _StructureViewState extends State<StructureView> {
  late List<Sp3dObj> objs = [];
  Sp3dWorld? world;
  Sp3dCamera? camera;
  bool isLoaded = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    api
        .downloadBinarySTL(widget.drugDetails.drugBankAccessionNumber)
        .then((stl) {
      if (!mounted) return;
      setState(() {
        final file = parseSTLData(stl);
        objs = [file.obj];
        world = Sp3dWorld(objs);
        camera = file.camera;
        isLoaded = true;
      });
    }).catchError((err) {
      if (!mounted) return;
      setState(() => _error = err);
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_error != null) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Text(
            _error!,
            style: Theme.of(context).textTheme.labelSmall,
          ),
        ),
      );
    } else if (!isLoaded) {
      return const Center(
        child: CircularProgressIndicator(),
      );
    } else {
      return LayoutBuilder(
        builder: (context, constraints) => Container(
          height: constraints.maxHeight,
          width: constraints.maxWidth,
          child: Sp3dRenderer(
            Size(constraints.maxWidth, constraints.maxHeight),
            Sp3dV2D(constraints.maxWidth / 2, constraints.maxHeight / 2),
            world!,
            Sp3dCamera(Sp3dV3D(0, 0, 3000), 6000),
            Sp3dLight(Sp3dV3D(0, 0, -1), syncCam: true),
          ),
        ),
      );
    }
  }
}
