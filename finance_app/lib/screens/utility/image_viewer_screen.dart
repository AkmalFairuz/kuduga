import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';
import 'package:photo_view/photo_view.dart';

class ImageViewerScreen extends StatelessWidget {
  const ImageViewerScreen(
      {Key? key, this.title = "Gambar", required this.provider})
      : super(key: key);

  final String title;
  final ImageProvider provider;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
          title: Text(title,
              style: AppBarTitle.style, overflow: TextOverflow.ellipsis)),
      body: Stack(
        children: [
          PhotoView(
            imageProvider: provider,
          )
        ],
      ),
    );
  }
}
