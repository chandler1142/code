CREATE TABLE `media`.`monitor_record` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `type` VARCHAR(45) NULL,
  `value` VARCHAR(100) NULL,
  `ip` VARCHAR(45) NULL,
  `create_time` DATETIME NULL,
  `properties` VARCHAR(200) NULL,
  PRIMARY KEY (`id`));
