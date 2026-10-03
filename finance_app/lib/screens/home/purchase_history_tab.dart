import 'package:finance_app/model/purchase.dart';
import 'package:finance_app/service/purchase_service.dart';
import 'package:finance_app/state/state_notifier.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/purchase/purchase_list.dart';
import 'package:finance_app/widgets/purchase/purchase_search.dart';
import 'package:finance_app/widgets/select/horizontal_select.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:flutter/material.dart';

import '../../widgets/text/text_app_bar.dart';

class PurchaseHistoryTab extends StatefulWidget {
  const PurchaseHistoryTab({super.key});

  @override
  State<PurchaseHistoryTab> createState() => _PurchaseHistoryTabState();
}

class _PurchaseHistoryTabState extends State<PurchaseHistoryTab> {
  List<PurchaseModel>? _purchases;
  int? _statusFilter;
  int? _lastId;
  final int _pageSize = 100;
  bool _isLast = false;
  late StateNotifierChannel _notifierChannel;

  @override
  void initState() {
    super.initState();
    _loadPurchases();

    _notifierChannel = StateNotifier.createChannel("purchase");
    _notifierChannel.addListener(_loadPurchases);
  }

  @override
  void dispose() {
    _notifierChannel.removeListener(_loadPurchases);
    super.dispose();
  }

  Future<void> _loadPurchases() async {
    _purchases = null;
    _lastId = null;
    _isLast = false;
    await _loadMorePurchases();
  }

  Future<void> _loadMorePurchases() async {
    var currentStatusFilter = _statusFilter;
    List<PurchaseModel> purchases = await PurchaseService.getPurchases(
        status: currentStatusFilter, limit: _pageSize, beforeId: _lastId);
    if (currentStatusFilter != _statusFilter) {
      return;
    }

    _isLast = purchases.length < _pageSize;
    _purchases ??= [];
    _purchases!.addAll(purchases);
    _lastId = purchases.lastOrNull?.purchaseId;
    setState(() {});
  }

  Future<void> _setStatusFilter(int? status) async {
    setState(() {
      _statusFilter = status;
      _purchases = null;
    });
    _loadPurchases();
  }

  Widget _buildSkeletonDate() {
    return Container(
      color: ColorHelper.bg100(context),
      padding: const EdgeInsets.all(16),
      child: const Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Skeleton(
            height: 12,
            width: 140,
          )
        ],
      ),
    );
  }

  Widget _buildSkeleton() {
    return Container(
      padding: const EdgeInsets.all(16),
      child: Row(
        children: [
          Skeleton(
            height: 35,
            width: 35,
            decoration: BoxDecoration(borderRadius: BorderRadius.circular(99)),
          ),
          const SizedBox(width: 16),
          const Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Skeleton(
                  height: 14,
                  width: 120,
                ),
                SizedBox(height: 8),
                Skeleton(
                  height: 8,
                  width: 80,
                ),
              ],
            ),
          ),
          const SizedBox(width: 16),
          const Skeleton(height: 14, width: 80)
        ],
      ),
    );
  }

  Widget _buildSkeletons() {
    return ListView(
      children: [
        _buildSkeletonDate(),
        const Line(),
        _buildSkeleton(),
        const Line(),
        _buildSkeleton(),
        const Line(),
        _buildSkeleton(),
        const Line(),
        _buildSkeleton(),
        const Line(),
        _buildSkeletonDate(),
        const Line(),
        _buildSkeleton(),
      ],
    );
  }

  void _onSearch() async {
    ModalBottomSheet.show(context, (context) => const PurchaseSearch());
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        AppBar(
          title: const AppBarTitle("Riwayat Pembelian"),
          actions: [
            IconButton(onPressed: _onSearch, icon: const Icon(Icons.search))
          ],
        ),
        const SizedBox(height: 4),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 8),
          child: HorizontalSelect(
              options: [
                HorizontalSelectOption(
                    value: null, widget: const Text("Semua")),
                HorizontalSelectOption(
                    value: 1, widget: const Text("Sedang diproses")),
                HorizontalSelectOption(value: 2, widget: const Text("Sukses")),
                HorizontalSelectOption(value: 3, widget: const Text("Gagal")),
              ],
              selected: _statusFilter,
              onSelect: (status) => _setStatusFilter(status)),
        ),
        const SizedBox(height: 4),
        const Line(),
        Expanded(
          child: RefreshIndicator(
              onRefresh: () async {
                setState(() {
                  _purchases = null;
                });
                await _loadPurchases();
              },
              child: _purchases != null
                  ? NotificationListener(
                      child: PurchaseList.build(context, _isLast, _purchases!),
                      onNotification: (notification) {
                        if (notification is ScrollEndNotification) {
                          if (_isLast) {
                            return false;
                          }
                          if (notification.metrics.pixels >=
                              notification.metrics.maxScrollExtent) {
                            _loadMorePurchases();
                          }
                        }
                        return false;
                      },
                    )
                  : _buildSkeletons()),
        )
      ],
    );
  }
}
