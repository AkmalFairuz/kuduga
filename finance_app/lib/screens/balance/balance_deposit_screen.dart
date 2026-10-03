import 'package:finance_app/screens/balance/deposit_info_screen.dart';
import 'package:finance_app/screens/balance/deposit_payment_screen.dart';
import 'package:finance_app/service/balance_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/currency_input.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/payment/payment_method.dart';
import 'package:finance_app/widgets/select/horizontal_select.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../../model/deposit.dart';
import '../../model/payment_method.dart';
import '../../utils/meta.dart';

class BalanceDepositScreen extends StatelessWidget {
  const BalanceDepositScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
        length: 2,
        child: Scaffold(
          resizeToAvoidBottomInset: false,
          appBar: AppBar(
            elevation: 1,
            title: const AppBarTitle('Deposit'),
            bottom: TabBar(
              onTap: (_) => Screens.unfocusInput(),
              labelStyle: Theme.of(context).textTheme.bodyLarge!.copyWith(
                  color: Colors.grey[50], fontWeight: FontWeight.w600),
              tabs: const [
                Tab(text: 'Isi Saldo'),
                Tab(text: 'Riwayat'),
              ],
            ),
          ),
          body: const TabBarView(
            children: [_DepositTab(), _DepositHistoryTab()],
          ),
        ));
  }
}

class _DepositTab extends StatefulWidget {
  const _DepositTab({Key? key}) : super(key: key);

  @override
  _DepositTabState createState() => _DepositTabState();
}

class _DepositTabState extends State<_DepositTab>
    with AutomaticKeepAliveClientMixin {
  List<PaymentMethodModel>? _paymentMethods;

  @override
  void initState() {
    BalanceService.getDepositPaymentMethods().then((value) {
      setState(() {
        _paymentMethods = value;
      });
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    if (_paymentMethods == null) {
      return const Center(
        child: CircularProgressIndicator(),
      );
    }
    return ListView.builder(
      itemBuilder: (context, index) {
        int index2 = index * 2;
        PaymentMethodModel paymentMethod1 = _paymentMethods![index2];
        PaymentMethodModel? paymentMethod2 =
            _paymentMethods!.length > (index2 + 1)
                ? _paymentMethods![index2 + 1]
                : null;
        return Container(
          padding: const EdgeInsets.only(left: 16, right: 16, top: 16),
          child: Row(
            children: [
              Expanded(
                child: PaymentMethod(
                  model: paymentMethod1,
                  onTap: () {
                    ModalBottomSheet.show(
                        context,
                        (context) => _DepositForm(
                              paymentMethod: paymentMethod1,
                            ));
                  },
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: paymentMethod2 != null
                    ? PaymentMethod(
                        model: paymentMethod2,
                        onTap: () {
                          ModalBottomSheet.show(
                              context,
                              (context) => _DepositForm(
                                    paymentMethod: paymentMethod2,
                                  ));
                        },
                      )
                    : Container(),
              ),
            ],
          ),
        );
      },
      itemCount: (_paymentMethods!.length / 2).ceil(),
    );
  }

  @override
  bool get wantKeepAlive => true;
}

class _DepositForm extends StatefulWidget {
  const _DepositForm({Key? key, required this.paymentMethod}) : super(key: key);

  final PaymentMethodModel paymentMethod;

  @override
  State<StatefulWidget> createState() => _DepositFormState();
}

class _DepositFormState extends State<_DepositForm> {
  final TextEditingController amountController = TextEditingController();

  void _onCreateDeposit() async {
    if (amountController.text.isEmpty) {
      Alert.message("Nominal deposit harus diisi");
      return;
    }
    final amount = int.parse(amountController.text);
    final paymentMethod = widget.paymentMethod;
    if (amount < paymentMethod.minAmount) {
      Alert.message(
          "Nominal deposit harus lebih dari ${Meta.currencyFormatRp(paymentMethod.minAmount)} (minimal nominal untuk metode pembayaran ${paymentMethod.name})");
      return;
    }
    if (amount > paymentMethod.maxAmount) {
      Alert.message(
          "Nominal deposit  harus kurang dari ${Meta.currencyFormatRp(paymentMethod.maxAmount)} (maksimal nominal untuk metode pembayaran ${paymentMethod.name})");
      return;
    }

    Alert.withLoading((done) async {
      DepositDetailedModel deposit =
          await BalanceService.createDeposit(amount, paymentMethod.id);
      done();
      Screens.back(); // close modal
      Screens.to(DepositPaymentScreen(deposit: deposit));
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SizedBox(height: 16),
        SizedBox(
          width: double.infinity,
          child: Column(
            children: [
              PaymentMethod.buildImage(widget.paymentMethod.imageUrl),
              const SizedBox(height: 12),
              Text("Deposit via ${widget.paymentMethod.name}")
            ],
          ),
        ),
        const SizedBox(height: 16),
        Text(
          "Min. Deposit: ${Meta.currencyFormatRp(widget.paymentMethod.minAmount)} | Maks. Deposit: ${Meta.currencyFormatRp(widget.paymentMethod.maxAmount)}",
          style: TextStyle(color: Colors.grey[500], fontSize: 12),
        ),
        const SizedBox(height: 8),
        CurrencyInput(
          hintText: "Masukkan nominal deposit",
          controller: amountController,
        ),
        const SizedBox(height: 16),
        Text(widget.paymentMethod.description,
            style: const TextStyle(fontSize: 12)),
        const SizedBox(height: 16),
        Row(
          children: [
            Expanded(
                child: Button(
                    variant: ButtonVariant.outline,
                    onPressed: () {
                      Screens.back();
                    },
                    child: const Text("Batal"))),
            const SizedBox(width: 8),
            Expanded(
                child: Button(
                    onPressed: _onCreateDeposit,
                    child: const Text("Lanjutkan")))
          ],
        )
      ],
    );
  }
}

class _DepositHistoryTab extends StatefulWidget {
  const _DepositHistoryTab({Key? key}) : super(key: key);

  @override
  _DepositHistoryTabState createState() => _DepositHistoryTabState();
}

class _DepositHistoryTabState extends State<_DepositHistoryTab>
    with AutomaticKeepAliveClientMixin {
  List<DepositModel>? _deposits;

  final List<_DepositHistoryFilterInfo> _filters = const [
    _DepositHistoryFilterInfo('Semua', null),
    _DepositHistoryFilterInfo('Menunggu', 0),
    _DepositHistoryFilterInfo('Berhasil', 1),
    _DepositHistoryFilterInfo('Gagal', 2),
  ];

  int? _statusFilter;

  @override
  void initState() {
    super.initState();
    _loadDeposit();
  }

  Future<void> _loadDeposit() {
    return BalanceService.getDepositHistory(status: _statusFilter)
        .then((value) {
      if (!mounted) {
        return;
      }
      setState(() => _deposits = value);
    }).catchError((err) {
      Alert.message(err.toString());
    });
  }

  Widget _buildSkeleton() {
    return Container(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        child: const Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Skeleton(width: 140, height: 17),
                  SizedBox(height: 11),
                  Skeleton(width: 70, height: 13),
                  SizedBox(height: 11),
                  Skeleton(width: 100, height: 13),
                ],
              ),
            ),
            Skeleton(
              width: 80,
              height: 15,
            )
          ],
        ));
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return Container(
        width: double.infinity,
        padding: const EdgeInsets.only(top: 4),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          const SizedBox(height: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12),
            child: HorizontalSelect(
              options: _filters
                  .map((filter) => HorizontalSelectOption(
                      value: filter.filter, widget: Text(filter.title)))
                  .toList(),
              selected: _statusFilter,
              onSelect: (filter) {
                setState(() {
                  _statusFilter = filter;
                  _deposits = null;
                });
                _loadDeposit();
              },
            ),
          ),
          const SizedBox(height: 8),
          Expanded(child: _depositList())
        ]));
  }

  Widget _depositList() {
    return RefreshIndicator(
        onRefresh: () async {
          setState(() {
            _deposits = null;
          });
          await _loadDeposit();
        },
        child: _deposits != null && _deposits!.isEmpty
            ? ListView(
                children: const [
                  SizedBox(height: 40),
                  Center(
                    child: Text("Tidak ada riwayat deposit yang ditemukan"),
                  ),
                ],
              )
            : ListView.separated(
                itemCount: _deposits != null ? _deposits!.length : 8,
                separatorBuilder: (_, __) => const Line(),
                itemBuilder: (context, index) {
                  if (_deposits == null) {
                    return _buildSkeleton();
                  }
                  DepositModel d = _deposits![index];
                  return _DepositHistoryItem(
                      onPress: () {
                        Alert.withLoading((done) async {
                          final depositDetailed =
                              await BalanceService.getDepositDetailed(d.id);
                          done();
                          if (depositDetailed.status == 0) {
                            Screens.to(
                                DepositPaymentScreen(deposit: depositDetailed));
                          } else {
                            Screens.to(
                                DepositInfoScreen(deposit: depositDetailed));
                          }
                        }, message: "Memuat deposit...");
                      },
                      status: d.status,
                      description: d.paymentMethod,
                      createdAt: d.createdAt,
                      amount: d.amount - d.fee);
                }));
  }

  @override
  bool get wantKeepAlive => true;
}

class _DepositHistoryFilterInfo {
  final String title;
  final int? filter;

  const _DepositHistoryFilterInfo(this.title, this.filter);
}

class _DepositHistoryItem extends StatelessWidget {
  const _DepositHistoryItem({
    Key? key,
    required this.onPress,
    required this.status,
    required this.description,
    required this.createdAt,
    required this.amount,
  }) : super(key: key);

  final void Function() onPress;
  final int status;
  final String description;
  final int createdAt;
  final int amount;

  @override
  Widget build(BuildContext context) {
    String statusText = DepositStatus.displayStatus(status);
    Color statusColor;
    IconData statusIcon;
    String dateFormat = DateFormat('dd/MM/yyyy HH:mm')
        .format(DateTime.fromMillisecondsSinceEpoch(createdAt * 1000));

    switch (status) {
      case 0:
        statusColor = Colors.amber[900]!;
        statusIcon = Icons.watch_later;
      case 1:
        statusColor = Colors.green;
        statusIcon = Icons.check_circle;
        break;
      case 2:
        statusColor = Colors.red;
        statusIcon = Icons.cancel;
        break;
      default:
        statusColor = Colors.grey;
        statusIcon = Icons.question_mark;
    }

    return InkWell(
      splashFactory: InkRipple.splashFactory,
      splashColor: ColorHelper.bg300(context).withOpacity(0.25),
      onTap: onPress,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        child: Row(
          children: [
            Expanded(
                child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // success badge
                Row(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Icon(
                      statusIcon,
                      color: statusColor,
                      size: 18,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      statusText,
                      style: const TextStyle(
                          fontSize: 14, fontWeight: FontWeight.bold),
                    ),
                  ],
                ),
                const SizedBox(height: 4),
                Text(
                  description,
                  style: const TextStyle(fontSize: 14),
                ),
                const SizedBox(height: 4),
                Text(
                  dateFormat,
                  style: const TextStyle(fontSize: 13, color: Colors.grey),
                ),
              ],
            )),
            Text(
              Meta.currencyFormatRp(amount),
              style: const TextStyle(fontSize: 14, fontWeight: FontWeight.bold),
            ),
          ],
        ),
      ),
    );
  }
}
