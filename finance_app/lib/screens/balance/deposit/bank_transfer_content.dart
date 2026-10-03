import 'package:finance_app/screens/balance/deposit/deposit_content.dart';
import 'package:flutter/material.dart';

import '../../../styles/color_helper.dart';
import '../../../utils/meta.dart';
import '../../../utils/toast.dart';
import '../../../widgets/text/text_link.dart';

class BankTransferContent extends DepositContent {
  const BankTransferContent(super.model, {super.key});

  @override
  State<BankTransferContent> createState() => _BankTransferContentState();
}

class _BankTransferContentState extends State<BankTransferContent> {
  @override
  Widget build(BuildContext context) {
    Color infoBg = ColorHelper.inputBg(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Rekening Bank Tujuan'),
        const SizedBox(height: 6),
        Text(widget.model.getBankAccountName() ?? '',
            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
        const SizedBox(height: 6),
        Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(8)),
            color: infoBg,
          ),
          child: Row(
            children: [
              Expanded(
                  child: Text(
                widget.model.getBankAccountNumber() ?? '',
                style:
                    const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
              )),
              // copy button
              TextLink(
                  onPressed: () {
                    // copy bank account number to clipboard
                    Meta.copyToClipboard(widget.model.getBankAccountNumber()!)
                        .then((_) => Toast.success(
                            context, 'Nomor rekening berhasil disalin'));
                  },
                  text: 'Salin'),
            ],
          ),
        ),
        const SizedBox(height: 16),
        const Text('Nominal Transfer (Rp)'),
        const SizedBox(height: 6),
        Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(8)),
            color: infoBg,
          ),
          child: Row(
            children: [
              Expanded(
                  child: Text(
                Meta.currencyFormat(widget.model.amount),
                style:
                    const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
              )),
              // copy button
              TextLink(
                  onPressed: () {
                    Meta.copyToClipboard(widget.model.amount.toString()).then(
                        (value) => Toast.success(
                            context, 'Nominal transfer berhasil disalin'));
                  },
                  text: 'Salin'),
            ],
          ),
        ),
        const SizedBox(height: 6),
        // alert that the nominal must be the same
        Row(
          children: [
            const SizedBox(width: 4),
            const Icon(Icons.warning, color: Colors.orange, size: 20),
            const SizedBox(width: 8),
            Expanded(
                child: Text(
              'Pastikan nominal transfer sama hingga 3 digit terakhir. Jumlah saldo yang diterima akan sama dengan nominal transfer. Pastikan transfer sebelum batas waktu',
              style: TextStyle(color: ColorHelper.bg400(context)),
            ))
          ],
        ),
      ],
    );
  }
}
