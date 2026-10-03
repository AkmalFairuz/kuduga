import 'package:flutter/material.dart';

abstract class Controller with ChangeNotifier {}

// Dyn is a wrapper for ListenableBuilder
// it's a shorthand for ListenableBuilder
class Dyn extends StatelessWidget {
  const Dyn(this.builder, this.listenable, {Key? key}) : super(key: key);

  final Widget? Function() builder;
  final ChangeNotifier listenable;

  @override
  Widget build(BuildContext context) {
    return ListenableBuilder(
        listenable: listenable,
        builder: (context, child) {
          Widget? returnWidget = builder();
          return returnWidget ?? const SizedBox();
        });
  }
}
