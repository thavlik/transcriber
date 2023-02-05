import 'package:demo/structure_view.dart';
import 'package:demo/testglb.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';
import 'home.dart';
import 'model.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return ScopedModel<MyModel>(
      model: MyModel(),
      child: ScopedModelDescendant(
        builder: (BuildContext context, Widget? child, MyModel model) =>
            MaterialApp(
          initialRoute: '/',
          title: 'Transcriber Demo',
          themeMode: ThemeMode.dark,
          darkTheme: ThemeData(
            brightness: Brightness.dark,
          ),
          theme: ThemeData(
            primarySwatch: Colors.blue,
          ),
          routes: {
            '/': (context) => const HomePage(),
          },
        ),
      ),
    );
  }
}
