import 'package:finance_app/model/product.dart';
import 'package:finance_app/screens/product/product_list_screen.dart';
import 'package:finance_app/screens/product/product_postpaid_list_screen.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/styles/text_style_helper.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/image/image.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import '../../widgets/layout/line.dart';

class ProductCategoryListScreen extends StatefulWidget {
  const ProductCategoryListScreen(
      {Key? key,
      this.isGrid = false,
      required this.title,
      required this.categories,
      this.fallbackImageUrl})
      : super(key: key);

  final String title;
  final bool isGrid;
  final List<ProductCategoryModel> categories;
  final String? fallbackImageUrl;

  @override
  State<ProductCategoryListScreen> createState() =>
      _ProductCategoryListScreenState();
}

class _ProductCategoryListScreenState extends State<ProductCategoryListScreen> {
  late List<ProductCategoryModel> categories;
  final TextEditingController searchController = TextEditingController();

  @override
  void initState() {
    categories = widget.categories;
    categories.sort((a, b) {
      return a.name.compareTo(b.name);
    });

    super.initState();
  }

  Widget _buildHeader() {
    return Container(
      color: Theme.of(context).appBarTheme.backgroundColor,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 16),
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.arrow_back, color: Colors.white),
                onPressed: () {
                  Screens.back();
                },
              ),
              const SizedBox(width: 16),
              Expanded(
                  child: Text(
                widget.title,
                style: TextStyleHelper.h1.copyWith(color: Colors.white),
              )),
              const SizedBox(width: 16),
              IconButton(
                icon: const Icon(Icons.home, color: Colors.white),
                onPressed: () =>
                    Navigator.of(context).popUntil((route) => route.isFirst),
              ),
            ],
          ),
          _buildSearch(),
        ],
      ),
    );
  }

  Widget _buildSearch() {
    return Container(
      padding: const EdgeInsets.all(16),
      child: Input(
        controller: searchController,
        onChanged: handleSearchChange,
        prefixIcon: const Icon(Icons.search),
        hintText: "Cari di ${widget.title}",
      ),
    );
  }

  void handleSearchChange(String value) {
    setState(() {
      categories = widget.categories
          .where((element) =>
              element.name.toLowerCase().contains(value.toLowerCase()))
          .toList();
    });
  }

  void onClickCategory(ProductCategoryModel category) {
    if (category.children.isEmpty) {
      if (category.isPostpaid()) {
        Screens.to(ProductPostpaidListScreen(category: category));
      } else {
        Screens.to(ProductListScreen(
          title: category.name,
          category: category,
          categoryImage:
              category.hasIcon() ? category.iconUrl : widget.fallbackImageUrl,
        ));
      }
      return;
    }
    if (category.children.isEmpty && !category.isPostpaid()) {
      Screens.to(ProductListScreen(
          title: category.name,
          category: category,
          categoryImage:
              category.hasIcon() ? category.iconUrl : widget.fallbackImageUrl));
      return;
    }
    Screens.to(ProductCategoryListScreen(
      isGrid: category.isGrid(),
      title: category.name,
      categories: category.children,
      fallbackImageUrl:
          category.hasIcon() ? category.iconUrl : widget.fallbackImageUrl,
    ));
  }

  Widget _buildContentGrid() {
    return Padding(
        padding: EdgeInsets.only(top: 8.sp, left: 8.sp, right: 8.sp),
        child: GridView.builder(
            itemCount: categories.length,
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 3, childAspectRatio: 0.8),
            itemBuilder: (context, index) {
              var category = categories[index];
              return GestureDetector(
                  onTap: () => onClickCategory(category),
                  child: Container(
                    padding: EdgeInsets.all(8.sp),
                    margin: EdgeInsets.all(6.sp),
                    decoration: BoxDecoration(
                        color: ColorHelper.bg0(context),
                        borderRadius: BorderRadius.circular(8.sp),
                        boxShadow: [
                          BoxShadow(
                              color:
                                  ColorHelper.bg200(context).withOpacity(0.3),
                              blurRadius: 8,
                              offset: const Offset(0, 4))
                        ],
                        border: Border.all(color: ColorHelper.bg200(context))),
                    child: Column(
                      children: [
                        XImage(
                          category.iconUrl,
                          width: 65.sp,
                          height: 65.sp,
                        ),
                        SizedBox(height: 19.sp),
                        Text(
                          category.name,
                          style: TextStyle(
                              fontSize: category.name.length >= 24
                                  ? 10.8.sp
                                  : 12.5.sp,
                              fontWeight: FontWeight.bold),
                          textAlign: TextAlign.center,
                        )
                      ],
                    ),
                  ));
            }));
  }

  Widget _buildContentNonGrid() {
    return ListView.separated(
        separatorBuilder: (_, __) => const Line(),
        itemCount: categories.length,
        itemBuilder: (context, index) {
          var category = categories[index];
          return ListTile(
            leading: XImage(
              category.hasIcon() ? category.iconUrl : widget.fallbackImageUrl!,
              width: 32,
              height: 32,
            ),
            title: Text(category.name),
            onTap: () => onClickCategory(category),
          );
        });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: SafeArea(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildHeader(),
          Expanded(
              child:
                  widget.isGrid ? _buildContentGrid() : _buildContentNonGrid())
        ],
      ),
    ));
  }
}
