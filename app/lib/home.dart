import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import 'model.dart';

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return ScopedModelDescendant<MyModel>(
      builder: (context, child, model) {
        return Scaffold(
          appBar: AppBar(
            title: Text(model.isConnected ? 'Connected' : 'Not Connected'),
          ),
          body: Container(),
        );
      },
    );
  }
}
