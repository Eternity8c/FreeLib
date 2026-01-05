import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';

@Component({
  selector: 'app-book-reader',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="reader-backdrop" (click)="onClose()">
      <div class="reader" (click)="$event.stopPropagation()">
        <div class="reader-header">
          <h3>{{ title }}</h3>
          <button class="close-btn" (click)="onClose()">✕</button>
        </div>
        <div class="reader-body">
          <ng-container *ngIf="loading">Загрузка...</ng-container>
          <pre *ngIf="textContent" class="reader-text">{{ textContent }}</pre>
          <iframe
            *ngIf="iframeSrc && !textContent"
            [src]="iframeSrc"
            class="reader-iframe"
          ></iframe>
          <div *ngIf="error" class="reader-error">
            Не удалось загрузить содержимое. <a (click)="openInNewTab()">Открыть в новой вкладке</a>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [
    `
      :host {
        position: fixed;
        inset: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 9999;
      }
      .reader-backdrop {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.4);
        display: flex;
        align-items: center;
        justify-content: center;
      }
      .reader {
        background: #fff;
        width: min(1000px, 95%);
        max-height: 90vh;
        border-radius: 8px;
        overflow: hidden;
        display: flex;
        flex-direction: column;
      }
      .reader-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0.75rem 1rem;
        border-bottom: 1px solid #eee;
      }
      .reader-body {
        padding: 1rem;
        overflow: auto;
        flex: 1;
      }
      .reader-iframe {
        width: 100%;
        height: 70vh;
        border: none;
      }
      .reader-text {
        white-space: pre-wrap;
        font-family: inherit;
      }
      .close-btn {
        background: transparent;
        border: none;
        font-size: 1.1rem;
        cursor: pointer;
      }
      .reader-error {
        color: #a00;
      }
    `,
  ],
})
export class BookReader implements OnInit {
  @Input() readUrl?: string | null;
  @Input() text?: string | null;
  @Input() title = 'Чтение';
  @Output() closed = new EventEmitter<void>();

  textContent: string | null = null;
  iframeSrc: SafeResourceUrl | null = null;
  loading = false;
  error = false;

  constructor(private sanitizer: DomSanitizer) {}

  ngOnInit(): void {
    // Если вызывающий передал текст напрямую, использовать его и пропустить загрузку
    if (this.text && this.text.length) {
      this.textContent = this.text;
      this.loading = false;
      return;
    }
    // В офлайн/локальном режиме не выполняем сетевые запросы.
    // Показать дружелюбное сообщение и дать возможность открыть readUrl в новой вкладке, если он доступен.
    if (!this.readUrl) {
      this.error = true;
      this.loading = false;
      return;
    }

    try {
      this.iframeSrc = this.sanitizer.bypassSecurityTrustResourceUrl(this.readUrl!);
    } catch (e) {
      this.error = true;
    }
    this.loading = false;
  }

  onClose() {
    this.closed.emit();
  }

  openInNewTab() {
    if (this.readUrl) window.open(this.readUrl, '_blank');
  }
}
