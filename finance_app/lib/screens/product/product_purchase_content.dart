import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/model/check.dart';
import 'package:finance_app/model/product.dart';
import 'package:finance_app/screens/auth/pin_screen.dart';
import 'package:finance_app/screens/purchase/purchase_detail_screen.dart';
import 'package:finance_app/service/purchase_service.dart';
import 'package:finance_app/state/state_notifier.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/fields/fields.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:flutter/material.dart';

import '../../utils/meta.dart';
import '../../widgets/button/button.dart';

class ProductPurchaseContent extends StatefulWidget {
  const ProductPurchaseContent(
      {Key? key,
      required this.product,
      required this.destination,
      required this.destinationFields})
      : super(key: key);

  final ProductModel product;
  final ProductDestinationModel destination;
  final List<ProductDestinationFieldRequestModel> destinationFields;

  @override
  State<ProductPurchaseContent> createState() => _ProductPurchaseContentState();
}

class _ProductPurchaseContentState extends State<ProductPurchaseContent> {
  // User can click "Buy" button after 1 second. This prevent accident click
  bool _canBuy = false;
  List<CheckResultModel>? checkResult;

  @override
  void initState() {
    super.initState();
    Future.delayed(const Duration(seconds: 1)).then((_) {
      if (mounted) {
        setState(() {
          _canBuy = true;
        });
      }
    });
    _checkDestination();
  }

  bool _isDestinationEmpty() {
    for (final field in widget.destinationFields) {
      if (field.value.isNotEmpty) {
        return false;
      }
    }
    return true;
  }

  void _checkDestination() async {
    if (widget.destination.checkerId == null || _isDestinationEmpty()) {
      return;
    }
    List<CheckResultModel> result = await PurchaseService.checkDestination(
        widget.destination.checkerId!, widget.destinationFields);
    setState(() {
      checkResult = result;
    });
  }

  void _onPurchase() async {
    String? pin;
    if (AccountController.getInstance().authDetails!.hasPin) {
      pin = await Screens.to(
          const PinScreen(title: "Masukkan PIN untuk melakukan pembelian"));
      if (pin == null) {
        return;
      }
    }

    try {
      Alert.withLoading((done) async {
        final purchaseId = await PurchaseService.createPurchase(
            widget.product.id, widget.destinationFields,
            pin: pin);
        done();
        StateNotifier.notify("purchase");
        Screens.back(); // close bottom modal
        Screens.to(PurchaseDetailScreen(purchaseId: purchaseId));
      });
    } finally {
      AccountController.getInstance().fetchAll();
    }
  }

  @override
  Widget build(BuildContext context) {
    ProductModel product = widget.product;
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Container(
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: const Text('Detail Pembelian',
            style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
      ),
      const SizedBox(height: 16),
      ...Fields.build(context, [
        FieldEntry("Produk", product.name),
        FieldEntry("Kategori", product.categoryName),
        FieldEntry("Harga", Meta.currencyFormatRp(product.price)),
        ...widget.destinationFields
            .map((e) => FieldEntry(e.label!, e.labelValue ?? e.value))
            .toList(),
        if (checkResult != null)
          ...checkResult!
              .where((e) => e.value.isNotEmpty)
              .map((e) => FieldEntry(e.name, e.value))
              .toList()
      ]),
      if (product.description.isNotEmpty) ...[
        const Line(),
        Container(
          padding: const EdgeInsets.all(16),
          width: double.infinity,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                product.description,
                style: TextStyle(
                    color: ColorHelper.theme(context,
                        dark: Colors.grey[300], light: Colors.grey[700])),
              )
            ],
          ),
        ),
        const Line(),
      ],
      const SizedBox(height: 16),
      Container(
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: Row(
          children: [
            Expanded(
                child: Button(
              onPressed: Navigator.of(context).pop,
              variant: ButtonVariant.outline,
              child: const Text("Batal"),
            )),
            const SizedBox(width: 8),
            Expanded(
                child: Button(
              isDisabled: !_canBuy,
              onPressed: _onPurchase,
              child: _canBuy
                  ? const Text("Beli")
                  : const Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                          Text(""),
                          SizedBox(
                              height: 16,
                              width: 16,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ))
                        ]),
            ))
          ],
        ),
      )
    ]);
  }
}
