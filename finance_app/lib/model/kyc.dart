class KycStatusModel {
  bool currentStatus;
  int? requestStatus;

  KycStatusModel({required this.currentStatus, this.requestStatus});

  factory KycStatusModel.fromJson(Map<String, dynamic> json) => KycStatusModel(
      currentStatus: json["currentStatus"],
      requestStatus: json["requestStatus"]);
}
