import 'dart:math';

import 'package:finance_app/utils/meta.dart';
import 'package:flutter/material.dart';

class XGrid extends StatelessWidget {
  const XGrid(
      {Key? key,
      required this.countPerRow,
      required this.children,
      this.horizontalGap = 0,
      this.verticalGap = 0})
      : super(key: key);

  final int countPerRow;
  final double horizontalGap;
  final double verticalGap;
  final List<Widget> children;

  @override
  Widget build(BuildContext context) {
    List<Widget> c = [];
    for (int i = 0; i < children.length; i += countPerRow) {
      final c2 = children.sublist(i, min(i + countPerRow, children.length));
      final w = Row(
          mainAxisAlignment: c2.length == countPerRow
              ? MainAxisAlignment.spaceBetween
              : MainAxisAlignment.center,
          children: c2
              .map((e) => Container(
                  padding: c2.length != countPerRow
                      ? EdgeInsets.symmetric(horizontal: horizontalGap)
                      : null,
                  child: e))
              .toList());
      c.add(w);
    }
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: Meta.addSeparatorToList(c, SizedBox(height: verticalGap)),
    );
  }
}
