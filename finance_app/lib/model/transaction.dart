class TransactionModel {
  final int id;
  final String description;
  final int amount;
  final int createdAt;
  final int beforeBalance;
  final int afterBalance;

  TransactionModel({
    required this.id,
    required this.description,
    required this.amount,
    required this.createdAt,
    required this.beforeBalance,
    required this.afterBalance,
  });

  bool isDebit() {
    return amount < 0;
  }

  bool isCredit() {
    return amount > 0;
  }

  factory TransactionModel.fromJson(Map<String, dynamic> json) {
    return TransactionModel(
      id: json['id'],
      description: json['description'],
      amount: json['amount'],
      createdAt: json['createdAt'],
      beforeBalance: json['beforeBalance'],
      afterBalance: json['afterBalance'],
    );
  }
}
