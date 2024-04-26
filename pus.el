;;; pus.el --- pus codeblock completion  -*- lexical-binding:t -*-

;;; Version: 1

;;; Commentary:

;; Enables evaluating org code blocks.
;;
;; Within the code blocks, enables completion at point that queries pus under the hood.

;;; Code:


(require 'org)
(require 'org-element)


(defcustom pus-path "pus"
  "Path to the pus executable."
  :group 'pus
  :type 'string)


(defcustom pus-default-directory "."
  "Default directory to execute pus in."
  :group 'pus
  :type 'string)


(defun pus-init ()
  "Overload `org-cycle' to do code completion in code blocks."
  (advice-add
   'org-cycle
   :around (lambda (orig-fun &rest args)
             (if (pus--in-source-block-p)
                 (pus-complete-at-point)
               (apply orig-fun args)))))


(defun org-babel-execute:pus (body _)
  "Execute BODY code with pus.
Ignore PARAMS."
  (let ((buf-inp (generate-new-buffer "*pus-input*" t))
        (buf-out (generate-new-buffer "*pus-output*" t))
        (buf-err (generate-new-buffer "*pus-error*" t))
        (result nil))
    (with-current-buffer buf-inp
      (insert body)
      (let ((default-directory pus-default-directory))
        (shell-command-on-region
         (point-min) (point-max)
         (format "%s --execute" pus-path)
         buf-out nil buf-err t)))
    (with-current-buffer buf-out
      (setq result (buffer-substring (point-min) (point-max))))
    (kill-buffer buf-inp)
    (kill-buffer buf-out)
    (kill-buffer buf-err)
    result))


(defun pus-complete-at-point ()
  "Complete pus code at point."
  (interactive)
  (let ((p (point))
        (cl (pus--current-text)))
    (pus--complete
     (plist-get cl ':text)
     (lambda (results)
       (let ((completion-at-point-functions
              (cons (lambda ()
                      (if (equal p (point))
                          (list
                           (save-excursion
                             (move-beginning-of-line nil)
                             (point))
                           (plist-get cl ':end)
                           (mapcar (lambda (x) x) results))
                        nil))
                    completion-at-point-functions)))
         (completion-at-point))))))


(defun pus--complete (code callback)
  "Pass CODE through pus completion and CALLBACK."
  (let ((default-directory pus-default-directory)
        (buf (generate-new-buffer "*pus-complete*" t))
        (proc nil))
    (setq proc (make-process
                :name "pus-process"
                :connection-type 'pipe
                :command (list pus-path "--complete" "-")
                :buffer buf
                :sentinel (lambda (_ _)
                            (let ((lines (json-parse-string
                                          (with-current-buffer buf
                                            (buffer-substring (point-min) (point-max))))))
                              (kill-buffer buf)
                              (funcall callback lines)))))
     (process-send-string proc code)
     (process-send-eof proc)))


(defun pus--current-text ()
  "Get current line or block text."
  (if (pus--in-source-block-p)
      (pus--current-block-up-to-point)
    (pus--current-line)))


(defun pus--current-line ()
  "Get current line text as a string.
Return a plist with :start and :end marking point bounds and
:text carrying the string with the current line text."
  (let ((p0 nil))
    (save-excursion
      (move-beginning-of-line nil)
      (setq p0 (point)))
    (list
     ':start p0
     ':end (point)
     ':text (buffer-substring-no-properties p0 (point)))))


(defun pus--current-block-up-to-point ()
  "Get code up to point in block like `pus--current-line'."
  (let ((block-start (caddr (cdddr (org-babel-get-src-block-info 'no-eval))))
        (block-code-start nil))
    (save-excursion
      (goto-char block-start)
      (forward-line)
      (setq block-code-start (point)))
    (list ':start block-code-start
          ':end (point)
          ':text (buffer-substring-no-properties block-code-start (point)))))


(defun pus--buffer-lines (b)
  "Convert a buffer B to a list of lines."
  (let ((lines (list)))
    (with-current-buffer b
      (goto-char (point-min))
      (while (not (eobp))
        (let ((line (buffer-substring-no-properties
                     (point)
                     (progn (forward-line 1) (point)))))
          (setq lines (cons line lines)))))
    (reverse lines)))


(defun pus--in-source-block-p ()
  "Check if point is inside a pus org-babel source code block."
  (let ((result nil))
    (org-babel-when-in-src-block
     (when (equal "pus" (car (org-babel-get-src-block-info 'no-eval)))
       (setq result t)))
    result))


(provide 'pus)
;;; pus.el ends here
