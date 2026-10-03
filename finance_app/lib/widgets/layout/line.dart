import 'package:finance_app/styles/color_helper.dart';
import 'package:flutter/material.dart';

class Line extends StatelessWidget {
  const Line({
    Key? key,
    this.height = 1,
    this.color,
  }) : super(key: key);

  final double height;
  final Color? color;

  @override
  Widget build(BuildContext context) {
    return Container(
        color: color ??
            ColorHelper.theme(context,
                dark: Colors.grey[700], light: Colors.grey[300])!,
        height: height,
        width: double.infinity);
  }
}
