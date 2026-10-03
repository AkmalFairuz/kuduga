import 'package:finance_app/model/stats.dart';
import 'package:finance_app/service/purchase_service.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

import '../../widgets/skeleton/skeleton.dart';

class StatsScreen extends StatelessWidget {
  const StatsScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: const TextAppBar("Statistik Anda"),
        body: SingleChildScrollView(
          child:
              Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            const SizedBox(height: 8),
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 16),
              child: Text("Statistik Pembelian",
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
            ),
            const SizedBox(height: 8),
            const Line(),
            FutureBuilder(
                future: PurchaseService.getStats(),
                builder: (context, snapshot) {
                  if (!snapshot.hasData) {
                    return Column(
                      children: [
                        _buildSkeletonContent(),
                        _buildSkeletonContent(),
                        _buildSkeletonContent(),
                        _buildSkeletonContent(),
                      ],
                    );
                  }
                  return Column(
                    children: [
                      _buildContent(context, "Hari Ini", snapshot.data!.today),
                      _buildContent(
                          context, "Kemarin", snapshot.data!.yesterday),
                      _buildContent(
                          context, "Bulan Ini", snapshot.data!.thisMonth),
                      _buildContent(context, "Bulan Sebelumnya",
                          snapshot.data!.previousMonth),
                    ],
                  );
                }),
          ]),
        ));
  }

  Widget _buildSkeletonContent() {
    return Container(
        margin: const EdgeInsets.all(16),
        child: const Skeleton(
          height: 180,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.all(Radius.circular(8)),
          ),
          width: double.infinity,
        ));
  }

  Widget _buildContent(
      BuildContext context, String title, PurchaseStatsModel model) {
    return Container(
        width: double.infinity,
        margin: const EdgeInsets.all(16),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(8),
          gradient: LinearGradient(
              begin: Alignment.bottomLeft,
              end: Alignment.topRight,
              colors: [
                Meta.color[400]!,
                Meta.color[800]!,
                Meta.color[900]!,
              ]),
        ),
        child: DefaultTextStyle(
          style: Theme.of(context).textTheme.bodyMedium!.copyWith(
                color: Colors.white,
              ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(title,
                  style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              const Line(),
              const SizedBox(height: 8),
              const Text("Total Keuntungan", style: TextStyle(fontSize: 14)),
              const SizedBox(height: 10),
              Text(Meta.currencyFormatRp(model.profit),
                  style: const TextStyle(
                      fontWeight: FontWeight.bold, fontSize: 24)),
              const SizedBox(height: 10),
              Text(
                  "${model.count} pembelian, total pembelian: ${Meta.currencyFormatRp(model.total)}")
            ],
          ),
        ));
  }
}
