import 'package:cached_network_image/cached_network_image.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:finance_app/service/app_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:flutter/material.dart';

class CarouselInfo extends StatelessWidget {
  const CarouselInfo({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 14),
      decoration: BoxDecoration(
          color: ColorHelper.bg100(context),
          border: Border.symmetric(
            horizontal: BorderSide(color: ColorHelper.bg200(context)),
          )),
      child: FutureBuilder(
          future: AppService.getBannerInfo(),
          builder: (context, snapshot) {
            return CarouselSlider(
              options: CarouselOptions(
                autoPlay: true,
                aspectRatio: 500 / 120,
                enlargeCenterPage: true,
                enlargeFactor: 0.2,
              ),
              items: snapshot.data == null
                  ? []
                  : snapshot.data!.map<Widget>((item) {
                      return _CarouselInfoItem(item.image);
                    }).toList(),
            );
          }),
    );
  }
}

class _CarouselInfoItem extends StatelessWidget {
  const _CarouselInfoItem(this.imageUrl);

  final String imageUrl;

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.all(1.0),
      child: ClipRRect(
          borderRadius: const BorderRadius.all(Radius.circular(5.0)),
          child: Stack(
            children: [
              Image(
                  image: CachedNetworkImageProvider(imageUrl),
                  fit: BoxFit.cover,
                  width: double.infinity),
            ],
          )),
    );
  }
}
