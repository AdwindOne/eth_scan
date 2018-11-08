/*
 Navicat Premium Data Transfer

 Source Server         : mac_localhost
 Source Server Type    : MySQL
 Source Server Version : 50723
 Source Host           : localhost
 Source Database       : eth

 Target Server Type    : MySQL
 Target Server Version : 50723
 File Encoding         : utf-8

 Date: 11/08/2018 09:57:43 AM
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `ec_address_log`
-- ----------------------------
DROP TABLE IF EXISTS `ec_address_log`;
CREATE TABLE `ec_address_log` (
  `state` varchar(1) DEFAULT NULL,
  `id` int(11) NOT NULL,
  `type` varchar(11) NOT NULL,
  `address` text,
  `pay_url` varchar(255) DEFAULT NULL,
  `is_delete` int(10) DEFAULT NULL,
  `pay_status` int(10) DEFAULT NULL,
  `last_block` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ec_crowd_order`
-- ----------------------------
DROP TABLE IF EXISTS `ec_crowd_order`;
CREATE TABLE `ec_crowd_order` (
  `id` int(11) NOT NULL,
  `order_type` varchar(255) DEFAULT NULL,
  `is_delete` varchar(255) DEFAULT NULL,
  `pay_url` varchar(255) DEFAULT NULL,
  `pay_status` varchar(255) DEFAULT NULL,
  `last_block` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `ec_ethdata`
-- ----------------------------
DROP TABLE IF EXISTS `ec_ethdata`;
CREATE TABLE `ec_ethdata` (
  `BlockNumber` int(11) DEFAULT NULL,
  `TimeStamp` int(11) DEFAULT NULL,
  `Hash` varchar(255) DEFAULT NULL,
  `Nonce` varchar(255) DEFAULT NULL,
  `BlockHash` varchar(255) DEFAULT NULL,
  `TransactionIndex` varchar(255) DEFAULT NULL,
  `From` varchar(255) DEFAULT NULL,
  `To` varchar(255) DEFAULT NULL,
  `Value` varchar(255) DEFAULT NULL,
  `Gas` varchar(255) DEFAULT NULL,
  `GasPrice` varchar(255) DEFAULT NULL,
  `Input` varchar(255) DEFAULT NULL,
  `ContractAddress` varchar(255) DEFAULT NULL,
  `CumulativeGasUsed` varchar(255) DEFAULT NULL,
  `GasUsed` varchar(255) DEFAULT NULL,
  `Confirmations` varchar(255) DEFAULT NULL,
  `IsError` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
