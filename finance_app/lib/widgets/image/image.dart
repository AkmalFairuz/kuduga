import 'package:cached_network_image/cached_network_image.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

class XImage extends StatelessWidget {
  const XImage(this.url, {Key? key, this.width, this.height}) : super(key: key);

  final double? width;
  final double? height;
  final String url;

  Widget _dummy(BuildContext context) {
    return SizedBox(
        width: width,
        height: height,
        child: Icon(
          CupertinoIcons.photo,
          size: width,
          color: Theme.of(context).colorScheme.primary,
        ));
  }

  @override
  Widget build(BuildContext context) {
    if (url.startsWith("http://") || url.startsWith("https://")) {
      return CachedNetworkImage(
          imageUrl: url,
          width: width,
          height: height,
          placeholder: (_, __) => Skeleton(
              width: width,
              height: height,
              decoration:
                  BoxDecoration(borderRadius: BorderRadius.circular(8))),
          errorWidget: (context, __, ___) => _dummy(context));
    }
    return _dummy(context);
  }
}
