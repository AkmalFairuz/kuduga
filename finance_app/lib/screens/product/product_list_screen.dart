import 'package:finance_app/model/product.dart';
import 'package:finance_app/screens/product/product_purchase_content.dart';
import 'package:finance_app/screens/purchase/purchase_postpaid_screen.dart';
import 'package:finance_app/service/product_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/bottom_select.dart';
import 'package:finance_app/utils/debouncer.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/utils/toast.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/product/product_small_tile.dart';
import 'package:finance_app/widgets/product/product_tile.dart';
import 'package:finance_app/widgets/select/horizontal_select.dart';
import 'package:flutter/material.dart';

import '../../model/contact.dart';
import '../../utils/alert.dart';
import '../../utils/screens.dart';
import '../../widgets/skeleton/skeleton.dart';
import '../../widgets/text/text_app_bar.dart';
import '../utility/select_contact_screen.dart';

class ProductListScreen extends StatefulWidget {
  static Future<ProductListScreen> show(int categoryId) async {
    var cat = await ProductService.getProductCategory(categoryId);
    return ProductListScreen(title: cat.name, category: cat);
  }

  const ProductListScreen(
      {Key? key,
      required this.title,
      required this.category,
      this.categoryImage,
      this.parentPulsaCategoryId,
      this.isPulsa = false})
      : super(key: key);

  final String title;
  final ProductCategoryModel? category;
  final int? parentPulsaCategoryId;
  final bool isPulsa;
  final String? categoryImage;

  @override
  State<ProductListScreen> createState() => _ProductListState();
}

class _ProductListState extends State<ProductListScreen> {
  static const sortByLowestPrice = 0;
  static const sortByHighestPrice = 1;
  static const sortByNameAZ = 2;

  String _selectedKind = "";
  int _selectedSortType = 0;

  ProductDestinationModel? _destination;
  final Map<String, ProductDestinationFieldModel> _destinationFields = {};
  final Map<String, TextEditingController> _destinationControllers = {};

  List<String>? _productKinds;
  List<ProductModel>? _products;
  ProductCategoryModel? _category;

  List<ProductCategoryModel>? _pulsaCategories;
  Map<String, String>? _pulsaPrefixes;
  int? _pulsaCategoryId;
  bool _pulsaLoading = false;
  int _parentPulsaCategoryId = 4;
  final int _pulsaDestinationId = 6;

  final _searchController = TextEditingController();
  final _searchFocusNode = FocusNode();
  bool _isSearching = false;
  final _searchDebounce = Debouncer(milliseconds: 100);
  List<ProductModel>? _filteredProducts;

  @override
  void initState() {
    if (widget.isPulsa) {
      _initPulsa();
      _parentPulsaCategoryId = widget.parentPulsaCategoryId ?? 4;
    } else {
      _category = widget.category!;
      _loadProducts();
    }
    super.initState();
  }

  void _initPulsa() async {
    final pulsaPrefixes = await ProductService.getPulsaCategories();
    final pulsaCategories = await ProductService.getProductCategories(
        parentId: _parentPulsaCategoryId);
    setState(() {
      _pulsaPrefixes = pulsaPrefixes;
      _pulsaCategories = pulsaCategories;
      _products = [];
    });
    _loadDestination(_pulsaDestinationId);
  }

  @override
  void dispose() {
    super.dispose();
    _searchDebounce.dispose();
  }

  Future<void> _loadDestination(int id) async {
    if (_destination != null && _destination!.id == id) {
      return;
    }
    var destination = await ProductService.getProductDestination(id);
    setState(() {
      _destination = destination;
      _destinationControllers.clear();
      _destinationFields.clear();

      for (final field in destination.fields) {
        _destinationFields[field.key] = field;
        _destinationControllers[field.key] = TextEditingController();
      }
    });
  }

  Future<void> _loadProducts() async {
    var att = 0;
    while (att < 8) {
      att++;
      try {
        _products = await ProductService.getProductByCategory(_category!.id);
        break;
      } catch (e) {
        if (mounted) {
          Toast.error(context, "Terjadi kesalahan dalam memuat produk: $e");
        }
      }
    }

    if (_products!.isNotEmpty) {
      _loadDestination(_products!.first.destinationId);
    }
    _productKinds = [];
    for (final product in _products!) {
      if (!_productKinds!.contains(product.kind)) {
        _productKinds!.add(product.kind);
      }
    }
    _selectedKind =
        _productKinds!.contains("") ? "" : (_productKinds!.firstOrNull ?? "");
    setState(() {});
  }

  void _resetPulsa() {
    setState(() {
      _pulsaCategoryId = null;
      _products = [];
      _category = null;
      _pulsaLoading = false;
    });
  }

  void _handleChangeDestination(String key, String newDestination) {
    if (!widget.isPulsa) {
      return;
    }
    if (newDestination.length < 4) {
      if (_pulsaCategoryId == null) {
        return;
      }
      _resetPulsa();
      return;
    }
    int? pulsaCategoryId = _detectPulsaCategory(newDestination);
    if (pulsaCategoryId != _pulsaCategoryId) {
      setState(() {
        _pulsaCategoryId = pulsaCategoryId;
        if (_pulsaCategoryId != null) {
          _pulsaLoading = true;
        } else {
          _products = [];
        }
      });
      if (pulsaCategoryId != null) {
        if (_pulsaCategoryId == pulsaCategoryId) {
          final cat = _pulsaCategories!
              .firstWhere((element) => element.id == pulsaCategoryId);
          _category = cat;
          setState(() {});
          _loadProducts().then((_) {
            if (_pulsaCategoryId != pulsaCategoryId) {
              setState(() {
                _resetPulsa();
              });
              return;
            }
            setState(() {
              _pulsaLoading = false;
            });
          });
        }
      } else {
        setState(() {
          _pulsaLoading = false;
          _category = null;
        });
      }
    }
  }

  int? _detectPulsaCategory(String phoneNumber) {
    for (var prefixes in _pulsaPrefixes!.entries) {
      if (phoneNumber.startsWith(prefixes.key)) {
        for (final cat in _pulsaCategories!) {
          if (cat.getMeta("pulsa") == prefixes.value) {
            return cat.id;
          }
        }
      }
    }
    return null;
  }

  List<ProductModel> getProducts() {
    if (_products == null) {
      return [];
    }
    return _products!.where((e) => e.kind == _selectedKind).toList();
  }

  List<Widget> _buildProductWidgets() {
    if (_products == null) {
      return [];
    }
    List<ProductModel> products = _filteredProducts ?? getProducts();
    switch (_selectedSortType) {
      case sortByLowestPrice:
        products = ProductModel.sortByPrice(products);
        break;
      case sortByHighestPrice:
        products = ProductModel.sortByPrice(products, lowest: false);
        break;
      case sortByNameAZ:
        products = ProductModel.sortByNameAZ(products);
        break;
    }
    if (widget.category != null && widget.category!.isGrid()) {
      return [
        Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            child: GridView.count(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              crossAxisCount: 2,
              mainAxisSpacing: 12,
              crossAxisSpacing: 12,
              childAspectRatio: 2.1,
              children: products.map((product) {
                return ProductSmallTile(
                    product: product,
                    onTap: () {
                      _onPrePurchase(context, product);
                    });
              }).toList(),
            ))
      ];
    }
    return products.map((product) {
      return Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
        child: ProductTile(
            product: product,
            fallbackImage: _category!.hasIcon()
                ? _category!.iconUrl
                : widget.categoryImage,
            onTap: () {
              _onPrePurchase(context, product);
            }),
      );
    }).toList();
  }

  Widget _buildProductSkeleton() {
    return Container(
        decoration: BoxDecoration(
            color: ColorHelper.theme(context,
                dark: Colors.grey[800], light: Colors.grey[100]),
            borderRadius: BorderRadius.circular(8)),
        padding: const EdgeInsets.all(16),
        margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
        child: Row(
          children: [
            Skeleton(
              width: 43,
              height: 43,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(16),
              ),
            ),
            const SizedBox(width: 16),
            const Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Skeleton(
                  width: 140,
                  height: 13,
                ),
                SizedBox(height: 8),
                Skeleton(
                  width: 60,
                  height: 10,
                ),
              ],
            )
          ],
        ));
  }

  void _onSearchChange(String query) {
    if (_products == null) {
      return;
    }
    _searchDebounce.run(() {
      query = query.toLowerCase();
      setState(() {
        _filteredProducts = getProducts()
            .where((element) => element.name.toLowerCase().contains(query))
            .toList();
      });
    });
  }

  Widget _buildSearchMenu() {
    return Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            const Icon(Icons.search),
            const SizedBox(width: 16),
            Expanded(
                child: TextField(
              focusNode: _searchFocusNode,
              controller: _searchController,
              onChanged: _onSearchChange,
              decoration: const InputDecoration(
                  hintText: "Cari produk",
                  isDense: true,
                  border: InputBorder.none,
                  contentPadding: EdgeInsets.zero),
            )),
          ],
        ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        resizeToAvoidBottomInset: false,
        appBar: AppBar(
          title: AppBarTitle(widget.title),
        ),
        body: ListView(children: [
          const SizedBox(height: 16),
          ..._buildDestinations(),
          _buildSelectMenu(),
          if (_isSearching) _buildSearchMenu(),
          const Line(),
          if (_productKinds != null)
            _productKinds!.length > 1
                ? Container(
                    padding: const EdgeInsets.only(left: 16, top: 8),
                    child: HorizontalSelect(
                        options: _productKinds!.map((e) {
                          return HorizontalSelectOption(
                              widget: Text(e == "" ? "Umum" : e), value: e);
                        }).toList(),
                        selected: _selectedKind,
                        onSelect: (v) {
                          if (_isSearching) {
                            _toggleSearch(false);
                          }
                          setState(() {
                            _selectedKind = v;
                          });
                        }),
                  )
                : const SizedBox(height: 8),
          if (_products == null ||
              (widget.isPulsa && (_pulsaPrefixes == null) ||
                  _pulsaLoading)) ...[
            const SizedBox(height: 4),
            _buildProductSkeleton(),
            _buildProductSkeleton(),
            _buildProductSkeleton(),
            _buildProductSkeleton(),
            _buildProductSkeleton(),
          ] else ...[
            ..._buildProductWidgets(),
          ],
        ]));
  }

  void _setDestination(String key, String destination) {
    _destinationControllers[key]!.text = destination;
    _handleChangeDestination(key, destination);
  }

  List<Widget> _buildDestinations() {
    if (_destination == null) {
      return [
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Skeleton(
            width: double.infinity,
            height: 48,
            decoration: BoxDecoration(borderRadius: BorderRadius.circular(4)),
          ),
        ),
        const SizedBox(height: 16)
      ];
    }
    return _destination!.fields.map((e) => _buildDestinationInput(e)).toList();
  }

  Widget _buildDestinationInput(ProductDestinationFieldModel field) {
    Widget? child;
    if (field.isOption()) {
      final bottomSelectItems = field.options!
          .map((e) => BottomSelectItem(value: e.value, widget: Text(e.label)))
          .toList();
      var valueToShow = "";
      for (final opt in field.options!) {
        if (opt.value == _destinationControllers[field.key]!.text) {
          valueToShow = opt.label;
          break;
        }
      }
      child = Input(
        hintText: field.label,
        enabled: false,
        controller: TextEditingController(text: valueToShow),
        suffixIcon: const Icon(Icons.arrow_drop_down_outlined, size: 24),
        onTap: () {
          BottomSelect.show(
                  context: context,
                  title: field.label,
                  selectedValue: _destinationControllers[field.key]!.text,
                  items: bottomSelectItems)
              .then((selected) {
            if (selected == null) {
              return;
            }
            setState(() {
              _destinationControllers[field.key]!.text = selected;
            });
          });
        },
      );
    } else {
      child = Input(
        hintText: field.label,
        keyboardType: field.getKeyboardType(),
        controller: _destinationControllers[field.key],
        onChanged: (text) => _handleChangeDestination(field.key, text),
        suffixIcon: field.isPhone()
            ? IconButton(
                splashRadius: 0.1,
                onPressed: () {
                  Screens.to<Contact>(const SelectContactScreen())
                      .then((contact) {
                    if (contact != null) {
                      _setDestination(field.key, contact.sanitizedPhone);
                    }
                  });
                },
                icon: const Icon(Icons.contacts),
              )
            : null,
      );
    }
    return Container(
        padding:
            const EdgeInsets.symmetric(horizontal: 16).copyWith(bottom: 16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            child,
            if (field.description.isNotEmpty) ...[
              const SizedBox(height: 4),
              Text(
                field.description,
                style:
                    TextStyle(fontSize: 12, color: ColorHelper.bg600(context)),
              )
            ]
          ],
        ));
  }

  Widget _buildSelectMenu() {
    return Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
        child: Row(
          children: [
            const Expanded(
              child: Text(
                "Pilih produk",
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 16,
                ),
              ),
            ),
            const SizedBox(width: 16),
            IconButton(
              padding: EdgeInsets.zero,
              constraints: const BoxConstraints(),
              splashRadius: 16,
              color: Meta.color[700]!,
              onPressed: () {
                BottomSelect.show(
                    context: context,
                    title: 'Urutkan',
                    selectedValue: _selectedSortType,
                    items: [
                      BottomSelectItem(
                          value: sortByLowestPrice,
                          widget: const Text('Harga paling rendah')),
                      BottomSelectItem(
                          value: sortByHighestPrice,
                          widget: const Text('Harga paling tinggi')),
                      BottomSelectItem(
                          value: sortByNameAZ, widget: const Text('Nama A-Z')),
                    ]).then((value) {
                  if (value != null) {
                    setState(() {
                      _selectedSortType = value;
                    });
                  }
                });
              },
              icon: const Icon(Icons.sort),
            ),
            if (_products != null && _products!.isNotEmpty) ...[
              const SizedBox(width: 16),
              IconButton(
                onPressed: () {
                  _toggleSearch(!_isSearching);
                },
                constraints: const BoxConstraints(),
                icon: Icon(_isSearching ? Icons.search_off : Icons.search),
                color: Meta.color[700]!,
                splashRadius: 16,
                padding: EdgeInsets.zero,
              )
            ]
          ],
        ));
  }

  void _toggleSearch(bool enable) {
    _searchController.clear();
    setState(() {
      _isSearching = enable;
      if (enable) {
        _filteredProducts = null;
      }
    });
    if (enable) {
      _searchFocusNode.requestFocus();
    }
  }

  void _onPrePurchase(BuildContext context, ProductModel product) {
    if (!product.isAvailable) {
      Alert.message(
          "Produk '${product.name}' tidak tersedia untuk sementara. Coba lagi nanti.");
      return;
    }
    if (_destination == null) {
      return;
    }
    List<ProductDestinationFieldRequestModel> destinationField = [];
    _destinationControllers.forEach((key, value) {
      ProductDestinationFieldModel dst = _destinationFields[key]!;
      String label = dst.label;
      String? labelValue;
      if (dst.isOption()) {
        for (final opt in dst.options!) {
          if (opt.value == value.text) {
            labelValue = opt.label;
          }
        }
      }
      destinationField.add(ProductDestinationFieldRequestModel(key, value.text,
          label: label, labelValue: labelValue));
    });
    if (product.isPostpaid()) {
      Screens.to(PurchasePostpaidScreen(
        productId: product.id,
        defaultDestination: _destinationControllers
            .map((key, value) => MapEntry(key, value.text)),
      ));
      return;
    }
    ModalBottomSheet.show(context,
        padding: const EdgeInsets.symmetric(vertical: 16), (context) {
      return ProductPurchaseContent(
        destination: _destination!,
        product: product,
        destinationFields: destinationField,
      );
    });
  }
}
