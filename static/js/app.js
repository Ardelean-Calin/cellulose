document.addEventListener('alpine:init', () => {
    Alpine.store('darkMode', {
        isDark: localStorage.getItem('darkMode') === 'true' || 
                (!localStorage.getItem('darkMode') && window.matchMedia('(prefers-color-scheme: dark)').matches),
        toggle() {
            this.isDark = !this.isDark;
            localStorage.setItem('darkMode', this.isDark);
            document.documentElement.classList.toggle('dark', this.isDark);
        },
        init() {
            document.documentElement.classList.toggle('dark', this.isDark);
            window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
                if (!localStorage.getItem('darkMode')) {
                    this.isDark = e.matches;
                    document.documentElement.classList.toggle('dark', e.matches);
                }
            });
        }
    });

    Alpine.data('tagsMenu', () => ({
        isOpen: false,
        searchTerm: '',
        tags: [],
        filteredTags: [],
        init() {
            this.tags = Array.from(this.$el.querySelectorAll('.tag-item')).map(el => ({
                name: el.querySelector('.tag-name').textContent,
                color: el.querySelector('.tag-color').classList[3].split('-')[1]
            }));
            this.filteredTags = this.tags;
        },
        searchTags() {
            const term = this.searchTerm.toLowerCase().trim();
            this.filteredTags = this.tags.filter(tag => 
                tag.name.toLowerCase().includes(term)
            );
        },
        hasExactMatch() {
            return this.searchTerm && this.filteredTags.some(tag => 
                tag.name.toLowerCase() === this.searchTerm.toLowerCase()
            );
        },
        close() {
            this.isOpen = false;
            this.searchTerm = '';
            this.filteredTags = this.tags;
        }
    }));

    Alpine.data('sorting', () => ({
        titleSort: null,
        dateSort: 'desc',
        sortByTitle() {
            this.titleSort = this.titleSort === 'asc' ? 'desc' : 'asc';
            this.dateSort = null;
            console.log('Sort by title:', this.titleSort);
        },
        sortByDate() {
            this.dateSort = this.dateSort === 'asc' ? 'desc' : 'asc';
            this.titleSort = null;
            console.log('Sort by date:', this.dateSort);
        }
    }));

    Alpine.data('dropzone', () => ({
        isDragging: false,
        uploadMessage: { text: '', type: 'success' },
        init() {
            window.addEventListener('keydown', (e) => {
                if (e.key === 'Escape' && this.isDragging) {
                    this.isDragging = false;
                }
            });
        },
        async handleDrop(event) {
            if (!this.isDragging) return;
            
            const file = event.dataTransfer.files[0];
            if (!file) return;

            const validTypes = ['application/pdf', 'image/jpeg', 'image/png', 'image/gif', 'image/webp'];
            const validExtensions = /\.(pdf|jpe?g|png|gif|webp)$/i;
            
            if (!validTypes.includes(file.type) && !file.name.match(validExtensions)) {
                this.showMessage('Only PDF and image files are supported', 'error');
                return;
            }

            const formData = new FormData();
            formData.append('pdf', file);

            try {
                const response = await fetch('/upload', {
                    method: 'POST',
                    body: formData
                });

                if (response.ok) {
                    this.showMessage(`File uploaded successfully (${this.formatSize(file.size)})`, 'success');
                    htmx.trigger(htmx.find("body"), "documentUploaded");
                } else {
                    throw new Error(await response.text());
                }
            } catch (err) {
                this.showMessage(err.message || 'Upload failed', 'error');
            }
        },
        showMessage(text, type = 'success') {
            this.uploadMessage = { text, type };
            setTimeout(() => this.uploadMessage = { text: '', type: 'success' }, 5000);
        },
        formatSize(bytes) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }
    }));
}); 