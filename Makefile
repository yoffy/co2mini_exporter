co2mini_exporter: *.go
	go build

.PHONY: clean
clean:
	$(RM) co2mini_exporter

.PHONY: install
install: co2mini_exporter
	if ! grep -q '^co2mini_exporter:' /etc/passwd; then \
		NOLOGIN_PATH="/usr/sbin/nologin"; \
		if [ ! -e "$${NOLOGIN_PATH}" ]; then \
			NOLOGIN_PATH="/sbin/nologin"; \
		fi; \
		useradd -N -M -d /var/run/co2mini_exporter -g plugdev -s "$${NOLOGIN_PATH}" co2mini_exporter; \
	fi
	install -m 755 co2mini_exporter /usr/local/bin/
	install -m 644 co2mini_exporter.service /etc/systemd/system/
	install -m 644 99-co2mini.rules /etc/udev/rules.d/
	udevadm control --reload-rules && udevadm trigger
	systemctl daemon-reload
	systemctl restart co2mini_exporter.service
	systemctl enable co2mini_exporter.service

.PHONY: uninstall
uninstall:
	systemctl disable co2mini_exporter.service
	systemctl stop co2mini_exporter.service
	$(RM) /usr/local/bin/co2mini_exporter
	$(RM) /etc/systemd/system/co2mini_exporter.service
	$(RM) /etc/udev/rules.d/99-co2mini.rules
	systemctl daemon-reload
	udevadm control --reload-rules && udevadm trigger
	userdel co2mini_exporter
