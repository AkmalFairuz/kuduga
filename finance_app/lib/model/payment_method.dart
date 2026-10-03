class PaymentMethodModel {
  int id;
  String name;
  String imageUrl;
  String description;
  bool noFee;
  int minAmount;
  int maxAmount;

  PaymentMethodModel(this.id, this.name, this.imageUrl, this.description,
      this.minAmount, this.maxAmount, this.noFee);

  factory PaymentMethodModel.fromJson(Map<String, dynamic> json) =>
      PaymentMethodModel(
          json["id"],
          json["name"],
          json["imageUrl"],
          json["description"],
          json["minAmount"],
          json["maxAmount"],
          json["noFee"]);
}
