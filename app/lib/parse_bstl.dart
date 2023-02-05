import 'dart:typed_data';
import 'package:simple_3d/simple_3d.dart';
import 'package:util_simple_3d/util_simple_3d.dart';
import 'package:simple_3d_renderer/simple_3d_renderer.dart';

// Parse binary STL data into a list of vertices and fragments
class STLFile {
  final Sp3dObj obj;
  final Sp3dCamera camera;

  STLFile(this.obj, this.camera);
}

STLFile parseSTLData(Uint8List stlData) {
  List<Sp3dV3D> vertices = [];
  List<Sp3dFragment> fragments = [];
  // The STL binary file format consists of a header, a number of
  // triangles, and the triangle data.
  // The header is 80 bytes long, and is ignored here.
  // The number of triangles is a 4-byte unsigned little endian integer
  // starting at byte 80.
  double minX = double.maxFinite,
      minY = double.maxFinite,
      maxX = -double.maxFinite,
      maxY = -double.maxFinite,
      minZ = double.maxFinite,
      maxZ = -double.maxFinite;

  int count = Uint32List.view(stlData.buffer, 80, 1).elementAt(0);
  for (int i = 0; i < count; i++) {
    final offset = i * 50 + 84;
    final vertBuf = ByteData.view(stlData.buffer, offset, 50);
    final index = vertices.length;

    final v1 = Sp3dV3D(
      vertBuf.getFloat32(12, Endian.little),
      vertBuf.getFloat32(16, Endian.little),
      vertBuf.getFloat32(20, Endian.little),
    );
    if (v1.x < minX) minX = v1.x;
    if (v1.y < minY) minY = v1.y;
    if (v1.x > maxX) maxX = v1.x;
    if (v1.y > maxY) maxY = v1.y;
    if (v1.z < minZ) minZ = v1.z;
    if (v1.z > maxZ) maxZ = v1.z;

    final v2 = Sp3dV3D(
      vertBuf.getFloat32(24, Endian.little),
      vertBuf.getFloat32(28, Endian.little),
      vertBuf.getFloat32(32, Endian.little),
    );
    if (v2.x < minX) minX = v2.x;
    if (v2.y < minY) minY = v2.y;
    if (v2.x > maxX) maxX = v2.x;
    if (v2.y > maxY) maxY = v2.y;
    if (v2.z < minZ) minZ = v2.z;
    if (v2.z > maxZ) maxZ = v2.z;

    final v3 = Sp3dV3D(
      vertBuf.getFloat32(36, Endian.little),
      vertBuf.getFloat32(40, Endian.little),
      vertBuf.getFloat32(44, Endian.little),
    );
    if (v3.x < minX) minX = v3.x;
    if (v3.y < minY) minY = v3.y;
    if (v3.x > maxX) maxX = v3.x;
    if (v3.y > maxY) maxY = v3.y;
    if (v3.z < minZ) minZ = v3.z;
    if (v3.z > maxZ) maxZ = v3.z;

    vertices.add(v1);
    vertices.add(v2);
    vertices.add(v3);
    fragments.add(Sp3dFragment([
      Sp3dFace([
        index + 0,
        index + 1,
        index + 2,
      ], 0)
    ]));
  }
  final obj = Sp3dObj(vertices, fragments, [
    FSp3dMaterial.grey.deepCopy(),
  ], []);
  final position = Sp3dV3D(
    minX + (maxX - minX) * 0.5,
    minY + (maxY - minY) * 0.5,
    minZ + (maxZ - minZ) * 0.5,
  );
  final camera = Sp3dCamera(position, 6000);
  return STLFile(obj, camera);
}
