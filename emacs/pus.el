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


(defun pus-complete-at-point ()
  "Complete pus code at point."
  (interactive)
  (let ((p (point))
        (cl (pus--current-line)))
    (pus--complete
     (plist-get cl ':text)
     (lambda (results)
       (let ((completion-at-point-functions
              (cons (lambda ()
                      (if (and p (point))
                          (list
                           (plist-get cl ':start)
                           (plist-get cl ':end)
                           results)
                        nil))
                    completion-at-point-functions)))
         (completion-at-point))))))


(defun pus--complete (code callback)
  "Pass CODE through pus completion and CALLBACK."
  (let ((default-directory pus-default-directory)
        (buf (generate-new-buffer "buf")))
    (make-process
     :name "pus-process"
     :command (list pus-path "--complete" code)
     :buffer buf
     :sentinel (lambda (x y)
                 (let ((lines (pus--buffer-lines buf)))
                   (kill-buffer buf)
                   (funcall callback (mapcar #'string-trim lines)))))))


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


;;; pus.el ends here
