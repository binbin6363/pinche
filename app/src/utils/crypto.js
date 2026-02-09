import CryptoJS from 'crypto-js'

/**
 * Hash password with MD5 before sending to server
 * @param {string} password - Plain text password
 * @returns {string} - MD5 hashed password
 */
export function hashPassword(password) {
  return CryptoJS.MD5(password).toString()
}
