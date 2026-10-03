import 'package:cached_network_image/cached_network_image.dart';
import 'package:finance_app/model/payment_method.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/widgets/image/cached_svg.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

class PaymentMethod extends StatelessWidget {
  const PaymentMethod({
    Key? key,
    required this.model,
    required this.onTap,
  }) : super(key: key);

  final PaymentMethodModel model;
  final VoidCallback onTap;

  static Widget buildImage(imageUrl) {
    Widget image = Container();
    switch (imageUrl.split(".").last) {
      case "svg":
        image = CachedSvg(imageUrl, height: 28.sp, width: double.infinity);
        break;
      case "png":
      case "jpg":
        image = CachedNetworkImage(
            imageUrl: imageUrl, height: 28.sp, width: double.infinity);
        break;
    }
    return image;
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      behavior: HitTestBehavior.translucent,
      onTap: onTap,
      child: Stack(
        children: [
          Container(
            height: 90.sp,
            width: double.infinity,
            padding: EdgeInsets.symmetric(vertical: 8.sp, horizontal: 16.sp),
            decoration: BoxDecoration(
              color: ColorHelper.theme(context,
                  dark: Colors.grey[800], light: Colors.grey[200]),
              borderRadius: BorderRadius.circular(8.sp),
            ),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                SizedBox(height: 4.sp),
                buildImage(model.imageUrl),
                SizedBox(height: 10.sp),
                Text(
                  model.name,
                  textAlign: TextAlign.center,
                  style: TextStyle(fontSize: 11.sp),
                )
              ],
            ),
          ),
          if (model.noFee)
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                Container(
                  padding:
                      EdgeInsets.symmetric(horizontal: 8.sp, vertical: 0.1.sp),
                  decoration: BoxDecoration(
                      color: Colors.green[600],
                      borderRadius: const BorderRadius.only(
                          topRight: Radius.circular(8),
                          bottomLeft: Radius.circular(8))),
                  child: Text(
                    'Bebas Admin',
                    style: TextStyle(
                        fontSize: 8.7.sp,
                        color: Colors.white,
                        fontWeight: FontWeight.bold),
                  ),
                )
              ],
            ),
        ],
      ),
    );
  }
}
