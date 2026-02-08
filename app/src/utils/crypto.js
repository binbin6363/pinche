import CryptoJS from 'crypto-js'

// salt for password hashing (should match backend)
const SALT = 'pinche_2024_salt'

/**
 * Hash password with SHA256 before sending to server
 * @param {string} password - Plain text password
 * @returns {string} - SHA256 hashed password
 */
export function hashPassword(password) {
  const salted = password + SALT
  return CryptoJS.SHA256(salted).toString(CryptoJS.enc.Hex)
}
