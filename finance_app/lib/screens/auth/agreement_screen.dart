import 'package:finance_app/service/app_service.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';

class AgreementScreen extends StatefulWidget {
  const AgreementScreen({Key? key, required this.onAgree}) : super(key: key);

  final VoidCallback onAgree;

  @override
  State<AgreementScreen> createState() => _AgreementScreenState();
}

class _AgreementScreenState extends State<AgreementScreen> {
  final webviewController = WebViewController();

  bool _loaded = false;
  String? _title;
  int _page = 0;

  @override
  void initState() {
    webviewController.setJavaScriptMode(JavaScriptMode.unrestricted);
    webviewController.setBackgroundColor(Colors.grey[50]!);
    webviewController.setNavigationDelegate(NavigationDelegate(
      onNavigationRequest: (request) {
        return NavigationDecision.prevent;
      },
      onPageFinished: (_) async {
        _title = (await webviewController.getTitle());
        _loaded = true;
        setState(() {});
        Meta.webviewInjectFont(webviewController);
      },
    ));
    loadWeb();
    super.initState();
  }

  Future<void> loadWeb() async {
    switch (_page) {
      case 0:
        final html = await AppService.getPrivacyPolicy();
        await webviewController.loadHtmlString(html);
        break;
      case 1:
        final html = await AppService.getTermsAndConditions();
        await webviewController.loadHtmlString(html);
        break;
    }
  }

  void onAgree() async {
    if (_page < 1) {
      _page++;
      _loaded = false;
      setState(() {});
      loadWeb();
      return;
    }
    widget.onAgree();
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
        child: Scaffold(
          appBar: TextAppBar(_title ?? "Persetujuan"),
          body: Column(
            children: [
              Expanded(
                  child: _loaded
                      ? WebViewWidget(
                          controller: webviewController,
                        )
                      : const Center(
                          child: CircularProgressIndicator(),
                        )),
              Container(
                padding: const EdgeInsets.all(16),
                child: Button(
                  width: double.infinity,
                  onPressed: onAgree,
                  isDisabled: !_loaded,
                  child: const Text("Saya setuju"),
                ),
              )
            ],
          ),
        ),
        onWillPop: () async {
          if (_page > 0) {
            _page--;
            await loadWeb();
            return false;
          }
          return true;
        });
  }
}
