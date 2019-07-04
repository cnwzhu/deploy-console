FROM  alpine
ADD console /root/console
RUN chmod o+x /root/console
CMD /root/console


