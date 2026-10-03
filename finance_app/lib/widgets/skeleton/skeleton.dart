import 'package:finance_app/styles/color_helper.dart';
import 'package:flutter/material.dart';
import 'package:shimmer/shimmer.dart';

class Skeleton extends StatelessWidget {
  const Skeleton({
    Key? key,
    this.height,
    this.width,
    this.baseColor,
    this.highlightColor,
    this.decoration = const BoxDecoration(),
  }) : super(key: key);

  final double? height;
  final double? width;
  final Color? baseColor;
  final Color? highlightColor;
  final BoxDecoration decoration;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: width,
      height: height,
      child: Shimmer.fromColors(
          baseColor: baseColor ??
              ColorHelper.theme(context,
                  dark: Colors.grey[600], light: Colors.grey[300])!,
          highlightColor: highlightColor ??
              ColorHelper.theme(context,
                  dark: Colors.grey[800], light: Colors.grey[100])!,
          child: Container(
            decoration: decoration.copyWith(
                color: baseColor ??
                    ColorHelper.theme(context,
                        dark: Colors.grey[600], light: Colors.grey[300])!),
          )),
    );
  }
}
