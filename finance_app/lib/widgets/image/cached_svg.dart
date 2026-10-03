import 'package:finance_app/utils/x_cache.dart';
import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

class CachedSvg extends StatefulWidget {
  const CachedSvg(this.url, {Key? key, this.width, this.height})
      : super(key: key);

  final double? width;
  final double? height;
  final String url;

  @override
  State<CachedSvg> createState() => _CachedSvgState();
}

class _CachedSvgState extends State<CachedSvg> {
  String data = "";

  @override
  void initState() {
    super.initState();
    XCache.network(widget.url).then((value) {
      setState(() {
        data = value;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    if (data == "") {
      return SizedBox(width: widget.width, height: widget.height);
    }
    return SvgPicture.string(data, width: widget.width, height: widget.height);
  }
}
