import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

class BrandIcon extends StatelessWidget {
  const BrandIcon({Key? key, this.size = 24}) : super(key: key);

  final double size;

  @override
  Widget build(BuildContext context) {
    return SvgPicture.asset("assets/icons/brand_white.svg",
        width: size, height: size);
  }
}
