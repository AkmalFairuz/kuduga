import 'dart:io';

import 'package:finance_app/model/purchase.dart';
import 'package:finance_app/screens/product/product_list_screen.dart';
import 'package:finance_app/screens/purchase/purchase_postpaid_screen.dart';
import 'package:finance_app/screens/support/create_ticket_screen.dart';
import 'package:finance_app/screens/utility/select_printer_screen.dart';
import 'package:finance_app/service/notification_service.dart';
import 'package:finance_app/service/purchase_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/bottom_buttons.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/fields/fields.dart';
import 'package:finance_app/widgets/input/currency_input.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:finance_app/widgets/text/modal_bottom_title.dart';
import 'package:flutter/material.dart';
import 'package:path_provider/path_provider.dart';
import 'package:screenshot/screenshot.dart';
import 'package:share_plus/share_plus.dart';

import '../../widgets/button/button.dart';

class PurchaseDetailScreen extends StatefulWidget {
  const PurchaseDetailScreen({Key? key, required this.purchaseId})
      : super(key: key);

  final int purchaseId;

  @override
  State<PurchaseDetailScreen> createState() => _PurchaseDetailState();
}

class _PurchaseDetailState extends State<PurchaseDetailScreen> {
  PurchaseModel? _purchase;
  ScreenshotController screenshotController = ScreenshotController();
  int lastScreenshot = 0;

  @override
  void initState() {
    _loadPurchase();

    NotificationService.instance.subscribeMessage(_handleStatusUpdate);

    super.initState();
  }

  @override
  void dispose() {
    NotificationService.instance.unsubscribeMessage(_handleStatusUpdate);
    super.dispose();
  }

  void _handleStatusUpdate(String event, Map<String, dynamic> data) {
    if (event == "purchaseStatusUpdate" && _purchase != null) {
      if (data["id"] as String == _purchase!.purchaseId.toString()) {
        _loadPurchase();
      }
    }
  }

  Future<void> _loadPurchase() async {
    var purchase = await PurchaseService.getPurchase(widget.purchaseId);
    if (!mounted) {
      return;
    }
    setState(() {
      _purchase = purchase;
    });
  }

  void _onShare() async {
    if (lastScreenshot + 1000 >= DateTime.now().millisecondsSinceEpoch) {
      return;
    }
    lastScreenshot = DateTime.now().millisecondsSinceEpoch;
    final Directory tempDir = await getTemporaryDirectory();
    var imagePath = await screenshotController.captureAndSave(tempDir.path,
        fileName: "purchase_${_purchase!.purchaseId.toString()}.png");
    if (imagePath == null) {
      return;
    }
    await Share.shareXFiles(
        [XFile(imagePath, name: "Purchase ${_purchase!.purchaseId}")]);
  }

  Widget _buildSkeletonField(double width) {
    return Skeleton(height: 14, width: width);
  }

  List<Widget> _buildBillInfo() {
    List<Widget> ret = [];
    if (_purchase!.billData == null) {
      return ret;
    }
    for (final data in _purchase!.billData!) {
      ret.addAll([
        Container(
          width: double.infinity,
          decoration: BoxDecoration(
              color: Theme.of(context).scaffoldBackgroundColor,
              border:
                  Border(top: BorderSide(color: ColorHelper.bg300(context))),
              borderRadius: const BorderRadius.only(
                  topLeft: Radius.circular(24), topRight: Radius.circular(24))),
          padding:
              const EdgeInsets.only(top: 12, left: 20, right: 20, bottom: 6),
          child: const Text("Detail Tagihan",
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 17)),
        ),
        ...Fields.build(context,
            data.entries.map((e) => FieldEntry(e.key, e.value)).toList()),
      ]);
    }
    return ret;
  }

  List<Widget> _buildNote() {
    List<Widget> ret = [];
    if (_purchase!.note == "") {
      return ret;
    }
    ret.addAll([
      const Line(),
      Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                "Catatan",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              Text(_purchase!.note),
            ],
          )),
      const Line(),
    ]);
    return ret;
  }

  List<Widget> _buildContent2() {
    List<FieldEntry> destinationInfo = [];
    if (_purchase != null) {
      _purchase!.destination.forEach((key, value) {
        destinationInfo.add(FieldEntry(key, value));
      });
    }
    return [
      SizedBox(
        height: 240,
        child: Stack(
          children: [
            Container(color: Meta.color[800]!, height: 240),
            Container(
              margin: const EdgeInsets.only(top: 220),
              decoration: BoxDecoration(
                  color: Theme.of(context).scaffoldBackgroundColor,
                  borderRadius: const BorderRadius.only(
                      topLeft: Radius.circular(20),
                      topRight: Radius.circular(20))),
              height: 20,
            ),
            Row(
              children: [
                IconButton(
                    onPressed: () {
                      Screens.back();
                    },
                    icon: const Icon(Icons.arrow_back, color: Colors.white)),
                const Expanded(child: SizedBox()),
                if (_purchase != null &&
                    _purchase!.status == PurchaseModel.statusSuccess) ...[
                  IconButton(
                      onPressed: () {
                        Screens.to(SelectPrinterScreen(printable: _purchase!));
                      },
                      icon: const Icon(Icons.print, color: Colors.white)),
                  IconButton(
                      onPressed: _onShare,
                      icon: const Icon(Icons.share, color: Colors.white))
                ],
              ],
            ),
            Center(
              child: IntrinsicHeight(
                child: Column(
                  children: [
                    _purchase != null
                        ? Container(
                            height: 75,
                            width: 75,
                            decoration: BoxDecoration(
                                color: Colors.white,
                                borderRadius: BorderRadius.circular(999)),
                            child: Icon(
                              _purchase!.statusIcon(),
                              color: Colors.blue[800]!,
                              size: 44,
                            ))
                        : Skeleton(
                            width: 75,
                            height: 75,
                            decoration: BoxDecoration(
                                borderRadius: BorderRadius.circular(999)),
                            highlightColor: Meta.color[400],
                            baseColor: Meta.color),
                    const SizedBox(height: 24),
                    _purchase != null
                        ? Text(_purchase!.statusText,
                            style: const TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 18))
                        : Skeleton(
                            height: 15,
                            width: 120,
                            highlightColor: Meta.color[400],
                            baseColor: Meta.color),
                    const SizedBox(height: 4),
                    _purchase != null
                        ? Text(
                            Meta.currencyFormatRp(_purchase!.isBill()
                                ? _purchase!.totalBill! +
                                    _purchase!.billAdmin! +
                                    _purchase!.userSellPrice
                                : _purchase!.userSellPrice),
                            style: const TextStyle(
                                color: Colors.white, fontSize: 15))
                        : Padding(
                            padding: const EdgeInsets.only(top: 8),
                            child: Skeleton(
                                height: 15,
                                width: 80,
                                highlightColor: Meta.color[400],
                                baseColor: Meta.color)),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
      ...Fields.build(context, [
        FieldEntry(
            'ID Pembelian',
            _purchase != null
                ? _purchase!.purchaseId.toString()
                : _buildSkeletonField(90)),
        FieldEntry(
            'Tanggal Pembelian',
            _purchase != null
                ? Meta.formatUnixDate(_purchase!.createdAt)
                : _buildSkeletonField(160)),
        FieldEntry(
            'Status',
            _purchase != null
                ? _purchase!.statusText
                : _buildSkeletonField(80)),
        FieldEntry(
            'Produk',
            _purchase != null
                ? _purchase!.productName
                : _buildSkeletonField(225)),
        FieldEntry(
            'Kategori Produk',
            _purchase != null
                ? _purchase!.productCategoryName
                : _buildSkeletonField(150)),
        if (_purchase == null)
          FieldEntry('Harga', _buildSkeletonField(100))
        else if (_purchase!.isBill()) ...[
          FieldEntry('Total Tagihan',
              Meta.currencyFormatRp(_purchase!.totalBill ?? 0)),
          FieldEntry(
              'Admin Bank', Meta.currencyFormatRp(_purchase!.billAdmin!)),
          FieldEntry(
              "Admin Loket",
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  Text(Meta.currencyFormatRp(_purchase!.userSellPrice)),
                  const SizedBox(width: 5),
                  GestureDetector(
                    onTap: _showPriceEditor,
                    child: const Icon(
                      Icons.edit,
                      size: 19,
                    ),
                  )
                ],
              )),
          FieldEntry(
              'Total Biaya',
              Meta.currencyFormatRp(_purchase!.totalBill! +
                  _purchase!.userSellPrice +
                  _purchase!.billAdmin!)),
        ] else
          FieldEntry(
              'Harga',
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  Text(Meta.currencyFormatRp(_purchase!.userSellPrice)),
                  const SizedBox(width: 5),
                  GestureDetector(
                    onTap: _showPriceEditor,
                    child: const Icon(
                      Icons.edit,
                      size: 19,
                    ),
                  )
                ],
              )),
        ...destinationInfo,
        if (_purchase != null && _purchase!.hasProof())
          FieldEntry('SN', _purchase!.proof!),
        if (_purchase != null && _purchase!.extraData.isNotEmpty)
          ..._purchase!.extraData.entries
              .map((e) => FieldEntry(e.key, e.value.toString()))
              .toList()
      ]),
      if (_purchase != null && _purchase!.isBill()) ..._buildBillInfo(),
      if (_purchase != null) ..._buildNote()
    ];
  }

  Widget _buildContent() {
    return Column(
      children: [
        Expanded(
            child: RefreshIndicator(
                onRefresh: () async {
                  setState(() {
                    _purchase = null;
                  });
                  await _loadPurchase();
                },
                child: SizedBox(
                  height: double.infinity,
                  child: SingleChildScrollView(
                      physics: const AlwaysScrollableScrollPhysics(),
                      child: Screenshot(
                        controller: screenshotController,
                        child: Container(
                            color: Theme.of(context).scaffoldBackgroundColor,
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: _buildContent2(),
                            )),
                      )),
                ))),
        Container(
          padding: const EdgeInsets.all(10),
          child: _purchase != null &&
                  _purchase!.status == PurchaseModel.statusSuccess
              ? Row(
                  children: [
                    Expanded(
                        child: Button(
                            onPressed: () async {
                              if (_purchase!.isBill()) {
                                Screens.to(PurchasePostpaidScreen(
                                  productId: _purchase!.productId,
                                ));
                              } else {
                                Alert.withLoading((done) async {
                                  var screen = await ProductListScreen.show(
                                      _purchase!.productCategoryId);
                                  done();
                                  Screens.to(screen);
                                });
                              }
                            },
                            child: const Text('Beli lagi'))),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Button(
                        variant: ButtonVariant.outline,
                        onPressed: _onNeedHelp,
                        child: const Text("Butuh Bantuan"),
                      ),
                    ),
                  ],
                )
              : null,
        ),
      ],
    );
  }

  void _onNeedHelp() {
    BottomButtons.show(context, [
      BottomButtonItem(
          text: "Status berhasil, tapi tidak masuk",
          onTap: () {
            Screens.to(CreateTicketScreen(
                initialProblemCategory: "purchase-invalid",
                initialMessage: "Purchase ID: ${_purchase!.purchaseId}\n"));
          }),
      BottomButtonItem(
          text: "Lainnya",
          onTap: () {
            Screens.to(const CreateTicketScreen());
          }),
    ]);
  }

  void _showPriceEditor() async {
    int? sellPrice = await ModalBottomSheet.show(context, (context) {
      return _UpdateSellPrice(purchase: _purchase!);
    }, showDragHelper: false);
    if (sellPrice == null) {
      return;
    }
    await PurchaseService.updateUserSellPrice(_purchase!.purchaseId, sellPrice);
    _purchase!.userSellPrice = sellPrice;
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: PreferredSize(
        preferredSize: const Size.fromHeight(0),
        child: AppBar(elevation: 0),
      ),
      body: _buildContent(),
    );
  }
}

class _UpdateSellPrice extends StatefulWidget {
  const _UpdateSellPrice({Key? key, required this.purchase}) : super(key: key);

  final PurchaseModel purchase;

  @override
  State<_UpdateSellPrice> createState() => _UpdateSellPriceState();
}

class _UpdateSellPriceState extends State<_UpdateSellPrice> {
  final TextEditingController _controller = TextEditingController();

  int _profit = 0;

  @override
  void initState() {
    super.initState();

    _controller.text = Meta.currencyFormat(widget.purchase.userSellPrice);
    _profit = _getProfit();
  }

  int _getProfit() {
    int curSellPrice = int.parse(_controller.text.replaceAll(".", ""));
    if (widget.purchase.isBill()) {
      return widget.purchase.billAdmin! +
          curSellPrice -
          widget.purchase.totalBillFee!;
    }
    return curSellPrice - widget.purchase.price;
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        ModalBottomTitle(
            widget.purchase.isBill() ? "Ubah Admin Loket" : "Ubah Harga Jual"),
        const SizedBox(height: 16),
        CurrencyInput(
          onChanged: (_) {
            setState(() {
              _profit = _getProfit();
            });
          },
          controller: _controller,
        ),
        const SizedBox(height: 16),
        if (widget.purchase.isBill())
          Text(
              "Admin bank ${Meta.currencyFormatRp(widget.purchase.billAdmin!)} dipotong ${Meta.currencyFormatRp(widget.purchase.billAdmin! - widget.purchase.totalBillFee!)} (komisi)\n"),
        Text(
            "${widget.purchase.isBill() ? "Jadi biaya admin yang anda bayar" : "Harga yang anda bayar"}: ${Meta.currencyFormatRp(widget.purchase.isBill() ? widget.purchase.totalBillFee! : widget.purchase.price)}"),
        if (widget.purchase.isBill())
          Text(
              "\nTotal yang anda bayar (total tagihan + admin bank - komisi): ${Meta.currencyFormatRp(widget.purchase.price)}\n"),
        Text("Keuntungan: ${Meta.currencyFormatRp(_profit)}",
            style: TextStyle(color: _profit > 0 ? Colors.green : Colors.red)),
        const SizedBox(height: 16),
        Button(
            width: double.infinity,
            onPressed: () {
              Navigator.of(context)
                  .pop(int.parse(_controller.text.replaceAll(".", "")));
            },
            child: const Text("Ubah")),
      ],
    );
  }
}
