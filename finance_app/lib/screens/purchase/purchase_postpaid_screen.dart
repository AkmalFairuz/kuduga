import 'package:finance_app/model/product.dart';
import 'package:finance_app/screens/purchase/purchase_detail_screen.dart';
import 'package:finance_app/screens/purchase/purchase_postpaid_confirm_screen.dart';
import 'package:finance_app/service/product_service.dart';
import 'package:finance_app/service/purchase_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/image/image.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

class PurchasePostpaidScreen extends StatefulWidget {
  const PurchasePostpaidScreen(
      {Key? key,
      required this.productId,
      this.defaultDestination,
      this.defaultDestinationId,
      this.product})
      : super(key: key);

  final int productId;
  final ProductModel? product;
  final Map<String, String>? defaultDestination;
  final int? defaultDestinationId;

  @override
  State<PurchasePostpaidScreen> createState() => _PurchasePostpaidState();
}

class _PurchasePostpaidState extends State<PurchasePostpaidScreen> {
  ProductModel? _product;
  ProductDestinationModel? _productDestination;
  final Map<String, TextEditingController> _destinationControllers = {};

  @override
  void initState() {
    super.initState();
    if (widget.product != null) {
      _product = widget.product;
      _loadProductDestination();
    } else {
      _loadProduct().then((_) {
        _loadProductDestination();
      });
    }
  }

  Future<void> _loadProduct() async {
    ProductModel product = await ProductService.getProduct(widget.productId);

    setState(() {
      _product = product;
    });
  }

  Future<void> _loadProductDestination() async {
    ProductDestinationModel productDestination =
        await ProductService.getProductDestination(_product!.destinationId);
    setState(() {
      _productDestination = productDestination;
      _destinationControllers.clear();
      for (var field in _productDestination!.fields) {
        String defaultText = "";
        if (widget.defaultDestination != null &&
            widget.defaultDestinationId == productDestination.id) {
          defaultText = widget.defaultDestination![field.key] ?? '';
        }
        _destinationControllers[field.key] =
            TextEditingController(text: defaultText);
      }
    });
  }

  Widget _buildProductInfo() {
    if (_product == null) {
      return const Row(
        children: [
          Skeleton(
            height: 40,
            width: 40,
          ),
          SizedBox(width: 16),
          Skeleton(
            height: 12,
            width: 120,
          )
        ],
      );
    }
    return Row(
      children: [
        XImage(_product!.imageUrl, width: 40, height: 40),
        const SizedBox(width: 16),
        Expanded(
            child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(_product!.name,
                style:
                    const TextStyle(fontWeight: FontWeight.w600, fontSize: 15)),
            Text("Biaya admin: ${Meta.currencyFormatRp(_product!.price)}",
                style: const TextStyle(fontSize: 11)),
          ],
        )),
      ],
    );
  }

  Widget _buildProductDestination(ProductDestinationFieldModel field) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Input(
        label: field.label,
        controller: _destinationControllers[field.key],
        keyboardType: field.getKeyboardType(),
      ),
    );
  }

  List<Widget> _buildProductDestinations() {
    if (_productDestination == null) {
      return [
        Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Skeleton(
                  height: 13,
                  width: 120,
                ),
                const SizedBox(height: 10),
                Skeleton(
                  height: 48,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(4),
                  ),
                  width: double.infinity,
                ),
              ],
            )),
      ];
    }
    return _productDestination!.fields
        .map((e) => _buildProductDestination(e))
        .toList();
  }

  void _handleContinue() async {
    List<ProductDestinationFieldRequestModel> destination = [];
    _destinationControllers.forEach((key, value) {
      destination.add(ProductDestinationFieldRequestModel(key, value.text));
    });

    Alert.withLoading((done) async {
      final bill = await PurchaseService.checkBill(_product!.id, destination);
      done();
      final purchaseId = await Screens.to(PurchasePostpaidConfirmScreen(
          bill: bill,
          product: _product!,
          productDestination: _productDestination!,
          destination: destination));
      if (purchaseId != null) {
        Screens.to(PurchaseDetailScreen(purchaseId: purchaseId));
      }
    }, message: "Mengecek tagihan...");
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      resizeToAvoidBottomInset: false,
      appBar: TextAppBar(_product == null ? "Tagihan" : _product!.categoryName),
      body: Column(
        children: [
          Expanded(
              child: ListView(
            children: [
              Container(
                padding: const EdgeInsets.all(16),
                child: _buildProductInfo(),
              ),
              const Line(),
              const SizedBox(height: 8),
              ..._buildProductDestinations(),
              const SizedBox(height: 8),
              const Line(),
            ],
          )),
          const Line(),
          Container(
            padding: const EdgeInsets.all(16),
            child: Button(
              onPressed: _handleContinue,
              width: double.infinity,
              isDisabled: _product == null,
              child: const Text("Lanjutkan"),
            ),
          )
        ],
      ),
    );
  }
}
