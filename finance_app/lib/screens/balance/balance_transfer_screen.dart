import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/controller/controller.dart';
import 'package:finance_app/screens/auth/pin_screen.dart';
import 'package:finance_app/screens/balance/balance_transfer_history_screen.dart';
import 'package:finance_app/screens/balance/balance_transfer_info_screen.dart';
import 'package:finance_app/service/balance_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/currency_input.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/select/select_box.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:finance_app/widgets/text/modal_bottom_title.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../model/transfer.dart';
import '../../utils/debouncer.dart';
import '../../widgets/input/input.dart';

class BalanceTransferScreen extends StatefulWidget {
  const BalanceTransferScreen({Key? key}) : super(key: key);

  @override
  State<StatefulWidget> createState() => _BalanceTransferState();
}

class _BalanceTransferState extends State<BalanceTransferScreen> {
  final queryController = TextEditingController();
  String oldQuery = "";
  final newDestinationController = TextEditingController();
  List<TransferUserInfoModel>? savedDestinations;
  List<TransferUserInfoModel>? filteredDestinations;
  final _debouncer = Debouncer(milliseconds: 250);

  @override
  void dispose() {
    super.dispose();
    _debouncer.dispose();
  }

  @override
  void initState() {
    super.initState();

    BalanceService.getTransferDestinations().then((dst) {
      setState(() {
        savedDestinations = dst;
      });
    });
  }

  List<TransferUserInfoModel> filteredSavedDestinations(String query) {
    if (query.isEmpty) {
      return savedDestinations!;
    }
    List<TransferUserInfoModel> lt;
    if (oldQuery.isNotEmpty &&
        filteredDestinations != null &&
        query.startsWith(oldQuery)) {
      lt = filteredDestinations!;
    } else {
      lt = savedDestinations!;
    }
    return lt
        .where((e) =>
            e.username.toLowerCase().contains(query) || e.email.contains(query))
        .toList();
  }

  Widget _buildList() {
    if (savedDestinations == null) {
      return const Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Skeleton(height: 18, width: 100),
                SizedBox(height: 12),
                Skeleton(height: 10, width: 150),
              ],
            ),
          ),
          Line(),
          Padding(
            padding: EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Skeleton(height: 18, width: 80),
                SizedBox(height: 12),
                Skeleton(height: 10, width: 130),
              ],
            ),
          ),
        ],
      );
    }
    var filtered = filteredDestinations ?? savedDestinations;
    return ListView.separated(
        itemBuilder: (context, index) {
          var s = filtered[index];
          return ListTile(
            onTap: () {
              Screens.to(_BalanceTransfer2Screen(destination: s));
            },
            title: Text(
              s.username,
              style: TextStyle(color: Meta.color[800]),
            ),
            subtitle: Text(s.email),
          );
        },
        separatorBuilder: (_, __) => const Line(),
        itemCount: filtered!.length);
  }

  void _onQueryChange(String query) {
    _debouncer.run(() => setState(() {
          query = query.toLowerCase();
          filteredDestinations = filteredSavedDestinations(query);
          oldQuery = query;
        }));
  }

  void _newDestination() {
    newDestinationController.clear();
    ModalBottomSheet.show(context,
        showDragHelper: false,
        isDismissible: false,
        enableDrag: false, (context) {
      return Column(
        children: [
          const ModalBottomTitle("Transfer ke tujuan baru", fontSize: 20),
          const SizedBox(height: 16),
          Input(
            controller: newDestinationController,
            hintText: "Alamat Email Tujuan",
          ),
          const SizedBox(height: 16),
          Button(
              onPressed: () async {
                String email = newDestinationController.text;
                if (!Meta.isValidEmail(email)) {
                  Alert.message("Email tidak valid");
                  return;
                }
                Alert.withLoading((done) async {
                  var destination =
                      await BalanceService.getTransferDestination(email);
                  done();
                  Screens.back(); // close modal
                  Screens.to(_BalanceTransfer2Screen(destination: destination));
                });
              },
              width: double.infinity,
              child: const Text("Lanjut"))
        ],
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      resizeToAvoidBottomInset: false,
      appBar: AppBar(
        title: const AppBarTitle("Transfer"),
        actions: [
          IconButton(
              onPressed: () {
                Screens.to(const BalanceTransferHistoryScreen());
              },
              icon: const Icon(Icons.history))
        ],
      ),
      body: Column(
        children: [
          const SizedBox(height: 16),
          Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Input(
                prefixIcon: const Icon(Icons.search),
                hintText: "Cari tujuan transfer",
                controller: queryController,
                onChanged: _onQueryChange,
              )),
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Button(
                onPressed: _newDestination,
                width: double.infinity,
                child: const Row(
                  children: [
                    Icon(Icons.person, color: Colors.white),
                    SizedBox(width: 16),
                    Expanded(
                        child: Text("Transfer ke tujuan baru",
                            textAlign: TextAlign.start)),
                    SizedBox(width: 16),
                    Icon(Icons.chevron_right, color: Colors.white)
                  ],
                )),
          ),
          const SizedBox(height: 16),
          const Line(),
          Expanded(
            child: _buildList(),
          ),
        ],
      ),
    );
  }
}

class _BalanceTransfer2Screen extends StatefulWidget {
  const _BalanceTransfer2Screen({Key? key, required this.destination})
      : super(key: key);

  final TransferUserInfoModel destination;

  @override
  State<StatefulWidget> createState() => _BalanceTransfer2State();
}

class _BalanceTransfer2State extends State<_BalanceTransfer2Screen> {
  final amountController = TextEditingController();
  final noteController = TextEditingController();

  void transferNow() async {
    if (amountController.text.isEmpty) {
      Alert.message("Nominal tidak boleh kosong");
      return;
    }
    int? amount = int.tryParse(amountController.text);
    if (amount == null || amount < 1) {
      Alert.message("Nominal tidak valid");
      return;
    }
    String note = noteController.text;

    String? pin;
    if (AccountController.getInstance().authDetails!.hasPin) {
      pin = await Screens.to(
          const PinScreen(title: "Masukkan PIN untuk melakukan transfer"));
      if (pin == null) {
        return;
      }
    }

    Alert.withLoading((done) async {
      var transfer = await BalanceService.transferToUser(
          destination: widget.destination,
          amount: amount,
          note: note,
          purpose: 1,
          pin: pin);
      done();
      AccountController.getInstance().fetchBalance();
      Screens.replace(BalanceTransferInfoScreen(model: transfer));
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        resizeToAvoidBottomInset: false,
        appBar: const TextAppBar("Transfer"),
        body: Column(
          children: [
            Expanded(
                child: ListView(
              children: [
                const SizedBox(height: 16),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text("Saldo anda saat ini:",
                          style: TextStyle(fontWeight: FontWeight.bold)),
                      Dyn(
                          () => Text(
                                Meta.currencyFormatRp(
                                    AccountController.getInstance().balance!),
                              ),
                          AccountController.getInstance())
                    ],
                  ),
                ),
                const SizedBox(height: 16),
                const Line(),
                const SizedBox(height: 16),
                const Padding(
                    padding: EdgeInsets.symmetric(horizontal: 16),
                    child: Text("Penerima:",
                        style: TextStyle(fontWeight: FontWeight.bold))),
                const SizedBox(height: 16),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  child: Row(
                    children: [
                      Icon(CupertinoIcons.person_alt_circle_fill,
                          size: 45, color: Meta.color),
                      const SizedBox(width: 10),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              widget.destination.username,
                              style: TextStyle(
                                  color: Meta.color[800],
                                  fontWeight: FontWeight.bold,
                                  fontSize: 16),
                            ),
                            const SizedBox(height: 2),
                            Text(
                              widget.destination.email,
                              style: TextStyle(color: Colors.grey[600]),
                            ),
                          ],
                        ),
                      )
                    ],
                  ),
                ),
                const SizedBox(height: 16),
                const Line(),
                const SizedBox(height: 16),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  child: Column(
                    children: [
                      CurrencyInput(
                        controller: amountController,
                        label: "Nominal",
                      ),
                      const SizedBox(height: 16),
                      Input(
                        controller: noteController,
                        label: "Catatan Transfer (opsional)",
                      ),
                      const SizedBox(height: 16),
                      const SelectBox(label: "Tujuan Transfer", items: [
                        SelectItem(1, "Transaksi"),
                        SelectItem(2, "Hadiah"),
                        SelectItem(3, "Keluarga"),
                        SelectItem(4, "Pembelian barang"),
                      ]),
                    ],
                  ),
                ),
              ],
            )),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Button(
                onPressed: transferNow,
                width: double.infinity,
                child: const Text("Kirim Sekarang"),
              ),
            )
          ],
        ));
  }
}
