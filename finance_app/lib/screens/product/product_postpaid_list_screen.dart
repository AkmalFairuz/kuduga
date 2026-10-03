import 'package:finance_app/model/product.dart';
import 'package:finance_app/screens/purchase/purchase_postpaid_screen.dart';
import 'package:finance_app/service/product_service.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/image/image.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

class ProductPostpaidListScreen extends StatefulWidget {
  const ProductPostpaidListScreen({Key? key, required this.category})
      : super(key: key);

  final ProductCategoryModel category;

  @override
  State<ProductPostpaidListScreen> createState() => _ProductPostpaidState();
}

class _ProductPostpaidState extends State<ProductPostpaidListScreen> {
  List<ProductModel>? _products;

  final searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _loadProducts();
  }

  Future<void> _loadProducts() async {
    var products =
        await ProductService.getProductByCategory(widget.category.id);
    products.sort((a, b) => a.name.compareTo(b.name));
    if (products.length == 1) {
      if (mounted) {
        Screens.replace(PurchasePostpaidScreen(
            productId: products.first.id, product: products.first));
        return;
      }
    }
    setState(() {
      _products = products;
    });
  }

  Widget _buildSkeleton() {
    return Container(
      padding: const EdgeInsets.all(16),
      child: Row(
        children: [
          Skeleton(
            width: 32,
            height: 32,
            decoration: BoxDecoration(borderRadius: BorderRadius.circular(4)),
          ),
          const SizedBox(width: 16),
          const Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Skeleton(width: 140, height: 16),
              SizedBox(height: 8),
              Skeleton(width: 80, height: 12),
            ],
          )
        ],
      ),
    );
  }

  Widget _buildContent() {
    if (_products == null) {
      return ListView(
        children: [
          _buildSkeleton(),
          const Line(),
          _buildSkeleton(),
          const Line(),
          _buildSkeleton(),
        ],
      );
    }
    final filteredProducts = searchController.text.isNotEmpty
        ? _products!
            .where((element) => element.name
                .toLowerCase()
                .contains(searchController.text.toLowerCase()))
            .toList()
        : _products!;
    if (filteredProducts.isEmpty) {
      return Center(
          child: Text("Produk '${searchController.text}' tidak ditemukan"));
    }
    return ListView.separated(
        itemBuilder: (context, index) {
          ProductModel product = filteredProducts[index];
          return ListTile(
            title: Text(product.name),
            onTap: () {
              Screens.to(PurchasePostpaidScreen(productId: product.id));
            },
            leading: XImage(product.imageUrl ?? '', width: 38, height: 38),
          );
        },
        separatorBuilder: (_, __) => const Line(),
        itemCount: filteredProducts.length);
  }

  void onSearchChange(String query) {
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: TextAppBar(widget.category.name),
      body: Column(
        children: [
          TextField(
            controller: searchController,
            autocorrect: false,
            onChanged: onSearchChange,
            decoration: InputDecoration(
                hintText: "Cari di ${widget.category.name}",
                prefixIcon: const Icon(Icons.search),
                suffixIcon: searchController.text.isNotEmpty
                    ? GestureDetector(
                        onTap: () {
                          searchController.clear();
                          setState(() {});
                        },
                        child: const Icon(Icons.close),
                      )
                    : null),
          ),
          Expanded(
            child: _buildContent(),
          )
        ],
      ),
    );
  }
}
