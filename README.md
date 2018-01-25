Для работы в контейнере надо модифицировать файл 600-upload.ini
добавив:

upload.virus_check = 1
general.av_scanner = "/var/www/rosneft/cometp/application/configs/config.d/avcheck /srv/etp/{filename}"

