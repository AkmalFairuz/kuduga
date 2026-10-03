import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/image/image.dart';
import 'package:flutter/material.dart';

import '../../model/product.dart';

class ProductTile extends StatelessWidget {
  const ProductTile({
    Key? key,
    required this.product,
    required this.onTap,
    this.fallbackImage,
  }) : super(key: key);

  final ProductModel product;
  final VoidCallback onTap;
  final String? fallbackImage;

  Widget _buildBadge({required Color color, required String text}) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 0.5),
      decoration: BoxDecoration(
          color: color,
          borderRadius: const BorderRadius.only(
            topRight: Radius.circular(8),
            bottomLeft: Radius.circular(8),
          )),
      child: Text(
        text,
        style: const TextStyle(
            fontSize: 10, fontWeight: FontWeight.bold, color: Colors.white),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    String? imageUrl = product.hasImage() ? product.imageUrl : fallbackImage;
    return Material(
        borderRadius: BorderRadius.circular(8),
        elevation: 0.4,
        shadowColor: Colors.grey[300],
        child: InkWell(
            borderRadius: BorderRadius.circular(8),
            onTap: onTap,
            child: Stack(
              alignment: Alignment.topRight,
              children: [
                if (!product.isAvailable)
                  _buildBadge(color: Colors.red[600]!, text: "Tidak tersedia"),
                Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(
                        color: ColorHelper.theme(context,
                            dark: Colors.grey[800], light: Colors.grey[200])!),
                  ),
                  padding:
                      const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                  child: Row(
                    children: [
                      imageUrl != null
                          ? XImage(imageUrl, width: 43, height: 43)
                          : Container(
                              width: 43,
                              height: 43,
                              decoration: BoxDecoration(
                                color: ColorHelper.theme(context,
                                    dark: Colors.grey[800]!,
                                    light: Colors.grey[100]!),
                                borderRadius: BorderRadius.circular(99),
                              ),
                              child: Icon(
                                Icons.question_mark,
                                color: Colors.grey[400],
                              )),
                      const SizedBox(width: 20),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(product.name,
                                style: const TextStyle(fontSize: 15)),
                            const SizedBox(height: 2),
                            Row(
                              crossAxisAlignment: CrossAxisAlignment.center,
                              children: [
                                Text(
                                  (product.isPostpaid() ? "+" : "") +
                                      Meta.currencyFormatRp(product.price),
                                  style: TextStyle(
                                    fontWeight: FontWeight.bold,
                                    fontSize: 13,
                                    color:
                                        Theme.of(context).colorScheme.primary,
                                  ),
                                ),
                                if (product.beforeDiscountPrice != null) ...[
                                  const SizedBox(width: 8),
                                  Stack(
                                    children: [
                                      Text(
                                        Meta.currencyFormatRp(
                                            product.beforeDiscountPrice!),
                                        style: TextStyle(
                                            fontSize: 11,
                                            color: ColorHelper.bg400(context),
                                            decoration:
                                                TextDecoration.lineThrough,
                                            decorationThickness: 3),
                                      ),
                                    ],
                                  )
                                ],
                              ],
                            )
                          ],
                        ),
                      ),
                    ],
                  ),
                )
              ],
            )));
  }
}
