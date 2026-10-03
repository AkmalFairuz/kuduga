import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:flutter/material.dart';

import '../../model/product.dart';

class ProductSmallTile extends StatelessWidget {
  const ProductSmallTile({Key? key, required this.product, required this.onTap})
      : super(key: key);

  final ProductModel product;
  final VoidCallback onTap;

  Widget _buildContent(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: ColorHelper.bg(context),
        boxShadow: [
          BoxShadow(
              color: ColorHelper.bg100(context),
              blurRadius: 8,
              offset: const Offset(0, 4))
        ],
        borderRadius: BorderRadius.circular(4),
        border: Border.all(color: ColorHelper.bg300(context)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
              child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Flexible(
                  child: Text(
                product.name,
                overflow: TextOverflow.fade,
                style: const TextStyle(fontWeight: FontWeight.bold),
              ))
            ],
          )),
          const SizedBox(height: 4),
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(Meta.currencyFormatRp(product.price),
                  textAlign: TextAlign.end,
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 13,
                    color: Theme.of(context).colorScheme.primary,
                  ))
            ],
          )
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
        onTap: onTap,
        child: Stack(
          alignment: Alignment.bottomLeft,
          children: [
            _buildContent(context),
            if (!product.isAvailable)
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 6, vertical: 0.5),
                decoration: BoxDecoration(
                  color: Colors.red[600],
                  borderRadius: const BorderRadius.only(
                      bottomLeft: Radius.circular(4),
                      topRight: Radius.circular(4)),
                ),
                child: const Text(
                  "Gangguan",
                  style: TextStyle(
                      color: Colors.white,
                      fontSize: 9.5,
                      fontWeight: FontWeight.bold),
                ),
              )
          ],
        ));
  }
}
