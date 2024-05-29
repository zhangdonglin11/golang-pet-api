/*
 Navicat Premium Data Transfer

 Source Server         : mySQL
 Source Server Type    : MySQL
 Source Server Version : 80035
 Source Host           : localhost:3306
 Source Schema         : golang_pet_db

 Target Server Type    : MySQL
 Target Server Version : 80035
 File Encoding         : 65001

 Date: 29/05/2024 18:59:17
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `pet_id` bigint UNSIGNED NULL DEFAULT NULL,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `target_id` bigint UNSIGNED NULL DEFAULT NULL,
  `level` tinyint NULL DEFAULT 0,
  `root_id` bigint UNSIGNED NULL DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `status` tinyint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_comment_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `fk_comment_target_user`(`target_id` ASC) USING BTREE,
  CONSTRAINT `fk_comment_target_user` FOREIGN KEY (`target_id`) REFERENCES `user` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of comment
-- ----------------------------
INSERT INTO `comment` VALUES (1, '2024-05-27 23:10:10.000', '2024-05-27 23:10:13.000', NULL, 1, 1, 2, 0, NULL, '111', 1);
INSERT INTO `comment` VALUES (2, '2024-05-27 23:10:53.000', '2024-05-27 23:10:56.000', NULL, 1, 2, 1, 1, 1, '222', 1);

-- ----------------------------
-- Table structure for pet
-- ----------------------------
DROP TABLE IF EXISTS `pet`;
CREATE TABLE `pet`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `pet_type` tinyint NULL DEFAULT NULL,
  `pet_breeds` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `pet_nickname` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `pet_gender` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `pet_age` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `pet_address` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `pet_status` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `pet_experience` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `pet_avatar` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `pet_intro` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `status` tinyint NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_pet_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pet
-- ----------------------------
INSERT INTO `pet` VALUES (1, '2024-05-25 19:35:03.000', NULL, NULL, 1, '哈士奇', '闪电哈', '弟弟', '一岁以内', '广东省·广州市·天河区', '找朋友', '无经验', NULL, '在家好无聊养只哈士奇', 2, 1);
INSERT INTO `pet` VALUES (2, '2024-05-25 20:17:10.086', '2024-05-25 22:19:56.382', NULL, 2, '中华田园猫', '咪咪', '妹妹', '一岁以内', '广东省·东莞市·茶山省', '找朋友', '没经验', '[\"http://localhost:8088/static/imageUtils/123.jpg\",\"http://localhost:8088/static/imageUtils/4456.png\"]', '', 2, 1);
INSERT INTO `pet` VALUES (4, '2024-05-25 22:24:18.745', '2024-05-25 22:24:18.745', NULL, 2, '中华田园猫', '喜喜', '妹妹', '一岁以内', '广东省·东莞市·中堂镇', '找对象', '没经验', '[\"http://localhost:8088/static/imageUtils/123.jpg\",\"http://localhost:8088/static/imageUtils/4456.png\"]', '', 2, 1);
INSERT INTO `pet` VALUES (5, '2024-05-25 22:24:51.711', '2024-05-25 22:24:51.711', NULL, 2, '中华田园猫', '喜喜', '妹妹', '一岁以内', '广东省·东莞市·中堂镇', '找对象', '没经验', '[\"http://localhost:8088/static/imageUtils/123.jpg\",\"http://localhost:8088/static/imageUtils/4456.png\"]', '', 1, 1);

-- ----------------------------
-- Table structure for pet_like
-- ----------------------------
DROP TABLE IF EXISTS `pet_like`;
CREATE TABLE `pet_like`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `pet_id` bigint UNSIGNED NULL DEFAULT NULL,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_pet_like_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pet_like
-- ----------------------------
INSERT INTO `pet_like` VALUES (1, '2024-05-26 23:42:10.000', '2024-05-26 23:42:13.000', NULL, 1, 1);
INSERT INTO `pet_like` VALUES (2, '2024-05-26 23:42:16.000', '2024-05-26 23:42:19.000', NULL, 2, 1);
INSERT INTO `pet_like` VALUES (3, '2024-05-26 23:42:32.000', '2024-05-26 23:42:34.000', NULL, 3, 1);

-- ----------------------------
-- Table structure for student
-- ----------------------------
DROP TABLE IF EXISTS `student`;
CREATE TABLE `student`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `age` bigint NULL DEFAULT NULL,
  `addr` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of student
-- ----------------------------
INSERT INTO `student` VALUES (1, '李四', 16, NULL);
INSERT INTO `student` VALUES (2, '王五', 18, '广东');
INSERT INTO `student` VALUES (3, '王五', 18, NULL);

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `nick_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `username` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `password` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `tel` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `gender` tinyint(1) NULL DEFAULT NULL,
  `role` bigint NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, '2024-04-27 22:47:34.991', '2024-04-27 22:47:34.991', NULL, '铲屎官7634', 'zdl7359', 'a123456', '1714229254988641900_6tHkKv8u.png', '', 0, 0, 1);
INSERT INTO `user` VALUES (2, '2024-05-27 23:17:04.000', '2024-05-27 23:17:07.000', NULL, '铲屎官1111', 'zdl1111', 'a123456', '1714229254988641900_6tHkKv8u.png', NULL, 1, 0, 1);

-- ----------------------------
-- Table structure for user_change
-- ----------------------------
DROP TABLE IF EXISTS `user_change`;
CREATE TABLE `user_change`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `nick_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `gender` bigint NULL DEFAULT NULL,
  `tel` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `approved` tinyint(1) NULL DEFAULT 1,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_user_user_change`(`user_id` ASC) USING BTREE,
  CONSTRAINT `fk_user_user_change` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_change
-- ----------------------------
INSERT INTO `user_change` VALUES (1, 1, '铲屎官111', '/static/image/1714321102526595600_W5eKcPUJ.png', 0, '13169197359', 1);

SET FOREIGN_KEY_CHECKS = 1;
