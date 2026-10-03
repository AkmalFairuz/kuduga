class PurchaseStatsWrapperModel {
  PurchaseStatsModel today;
  PurchaseStatsModel yesterday;
  PurchaseStatsModel thisMonth;
  PurchaseStatsModel previousMonth;

  PurchaseStatsWrapperModel(
      {required this.today,
      required this.yesterday,
      required this.thisMonth,
      required this.previousMonth});

  factory PurchaseStatsWrapperModel.fromJson(Map<String, dynamic> json) =>
      PurchaseStatsWrapperModel(
          today: PurchaseStatsModel.fromJson(json["today"]),
          yesterday: PurchaseStatsModel.fromJson(json["yesterday"]),
          thisMonth: PurchaseStatsModel.fromJson(json["thisMonth"]),
          previousMonth: PurchaseStatsModel.fromJson(json["previousMonth"]));
}

class PurchaseStatsModel {
  int total;
  int count;
  int profit;

  PurchaseStatsModel(
      {required this.total, required this.count, required this.profit});

  factory PurchaseStatsModel.fromJson(Map<String, dynamic> json) =>
      PurchaseStatsModel(
          total: json["total"], count: json["count"], profit: json["profit"]);
}
