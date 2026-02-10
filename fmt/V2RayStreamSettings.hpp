#pragma once

#include "AbstractBean.hpp"

namespace NekoGui_fmt {
    class V2rayStreamSettings : public JsonStore {
    public:
        QString network = "tcp";
        QString security = "";
        QString packet_encoding = "";
        // ws/http/grpc/tcp-http/httpupgrade
        QString path = "";
        QString host = "";
        // kcp/quic/tcp-http
        QString header_type = "";
        // tls
        QString sni = "";
        QString alpn = "";
        QString certificate = "";
        QString utlsFingerprint = "";
        bool allow_insecure = false;
        // ws early data
        QString ws_early_data_name = "";
        int ws_early_data_length = 0;
        // reality
        QString reality_pbk = "";
        QString reality_sid = "";
        QString reality_spx = "";
        // ech
        bool ech_enabled = false;
        QString ech_config = "";
        QString ech_config_path = "";
        QString ech_query_server_name = "";
        QString ech_bootstrap_resolver = "dns-direct";
        QString ech_tunnel_resolver = "dns-remote";
        // tls advanced
        bool disable_sni = false;
        QString tls_min_version = "";
        QString tls_max_version = "";
        // dialer options
        bool tcp_fast_open = false;
        bool udp_fragment = false;
        // multiplex
        int multiplex_status = 0;

        V2rayStreamSettings() : JsonStore() {
            _add(new configItem("net", &network, itemType::string));
            _add(new configItem("sec", &security, itemType::string));
            _add(new configItem("pac_enc", &packet_encoding, itemType::string));
            _add(new configItem("path", &path, itemType::string));
            _add(new configItem("host", &host, itemType::string));
            _add(new configItem("sni", &sni, itemType::string));
            _add(new configItem("alpn", &alpn, itemType::string));
            _add(new configItem("cert", &certificate, itemType::string));
            _add(new configItem("insecure", &allow_insecure, itemType::boolean));
            _add(new configItem("h_type", &header_type, itemType::string));
            _add(new configItem("ed_name", &ws_early_data_name, itemType::string));
            _add(new configItem("ed_len", &ws_early_data_length, itemType::integer));
            _add(new configItem("utls", &utlsFingerprint, itemType::string));
            _add(new configItem("pbk", &reality_pbk, itemType::string));
            _add(new configItem("sid", &reality_sid, itemType::string));
            _add(new configItem("spx", &reality_spx, itemType::string));
            _add(new configItem("ech_en", &ech_enabled, itemType::boolean));
            _add(new configItem("ech_cfg", &ech_config, itemType::string));
            _add(new configItem("ech_cfg_path", &ech_config_path, itemType::string));
            _add(new configItem("ech_query", &ech_query_server_name, itemType::string));
            _add(new configItem("ech_boot", &ech_bootstrap_resolver, itemType::string));
            _add(new configItem("ech_tun", &ech_tunnel_resolver, itemType::string));
            _add(new configItem("disable_sni", &disable_sni, itemType::boolean));
            _add(new configItem("tls_min_ver", &tls_min_version, itemType::string));
            _add(new configItem("tls_max_ver", &tls_max_version, itemType::string));
            _add(new configItem("tcp_fast_open", &tcp_fast_open, itemType::boolean));
            _add(new configItem("udp_fragment", &udp_fragment, itemType::boolean));
            _add(new configItem("mux_s", &multiplex_status, itemType::integer));
        }

        void BuildStreamSettingsSingBox(QJsonObject *outbound);
    };

    inline V2rayStreamSettings *GetStreamSettings(AbstractBean *bean) {
        if (bean == nullptr) return nullptr;
        auto stream_item = bean->_get("stream");
        if (stream_item != nullptr) {
            auto stream_store = (JsonStore *) stream_item->ptr;
            auto stream = (NekoGui_fmt::V2rayStreamSettings *) stream_store;
            return stream;
        }
        return nullptr;
    }
} // namespace NekoGui_fmt
