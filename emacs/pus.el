;;; pus.el -- interfacing with pus from within emacs  -*- lexical-binding:t -*-

;;; Commentary:

;; Enables evaluating org code blocks.
;;
;; Within the code blocks, enables completion at point that queries pus under the hood.

;;; Code:


(defcustom pus-path "pus"
  "Path to the pus executable.")


(defcustom pus-default-directory "."
  "Default directory to execute pus in.")


(defun org-babel-execute:pus (body params)
  "Execute BODY code with pus.
Ignore PARAMS."
  (let ((default-directory pus-default-directory))
    (shell-command-to-string
     (format "%s --execute %s" pus-path
             (shell-quote-argument body)))))


;;; pus.el ends here
