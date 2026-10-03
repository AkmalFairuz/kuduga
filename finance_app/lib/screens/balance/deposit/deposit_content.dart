import 'package:finance_app/model/deposit.dart';
import 'package:flutter/material.dart';

abstract class DepositContent extends StatefulWidget {
  const DepositContent(this.model, {Key? key}) : super(key: key);

  final DepositDetailedModel model;
}
